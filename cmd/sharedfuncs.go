package cmd

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		Log.Errorf("Error getting interfaces: %v\n", err)
		return "", fmt.Errorf("error getting interfaces: %v", err)
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			Log.Errorf("Error getting addresses: %v\n", err)
			return "", fmt.Errorf("error getting addresses: %v", err)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Check if the IP address is not a loopback address
			if ip != nil && !ip.IsLoopback() {
				if ip.To4() != nil {
					return ip.String(), nil
				}
			}
		}
	}
	Log.Error("No IP address found...")
	return "", fmt.Errorf("no IP address found")
}

func CheckHttpServerPort() string {
	// Description: Check if the HTTP server port is open and responds with 200 OK
	HttpServerFull = fmt.Sprintf("http://%s:%d", KaliMachine, HttpServerPort)
	Log.Debugf("Checking status on: %s\n", HttpServerFull)
	resp, err := http.Get(HttpServerFull)
	if err != nil {
		Log.Errorf("Error checking HTTP server: %v\n", err)
		return "Error"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		Log.Debugf("Status OK on: %s\n", HttpServerFull)
		return "OK"
	} else {
		return string(resp.StatusCode)
	}
}

func CheckHTTPServerType() (bool, error) {
	// Description: Check if the HTTP server is running on Windows or Linux
	var url string
	if UploadServerPort > 0 {
		url = fmt.Sprintf("http://%s:%d/", KaliMachine, UploadServerPort)
	} else {
		url = fmt.Sprintf("http://%s:%d/Workspace", KaliMachine, HttpServerPort)
	}
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		Log.Errorf("Error making POST request: %v\n", err)
		return false, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return true, nil
	} else {
		return false, nil
	}
}

func AskForInput() {

	// Terminate if no Connection
	if UploadServerStatus == "Error" && CheckHttpServerPort() == "Error" {
		Log.Fatal("No connection with Web Servers. Terminating...")
		os.Exit(1)
	}

	CheckWorkspaceDirectories()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to run PEAS? (y/N): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "yes" || input == "y" {
			RunPEAS()
		} else if input == "no" || input == "N" || input == "n" {
			Log.Printf("Downloading all Tools...")
			DownloadAllTools()
			os.Exit(0)
		} else {
			fmt.Println("Please type 'yes' or 'no' and press enter.")
		}
	}
}

func UnzipDirToPath(zipFile string, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		Log.Debugf("Found file: %s\n", f.Name)
		if f.Name == "" || f.Name == "./" {
			continue
		}
		fpath := filepath.Join(destDir, f.Name)
		Log.Debugf("Extracting file: %s\n", fpath)
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			Log.Errorf("Illegal file path: %s\n", fpath)
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteZipFiles(zipFile string) error {
	// Delete the old zip file after successful extraction
	err := os.Remove(zipFile)
	if err != nil {
		Log.Errorf("Failed to delete zip file: %v\n", err)
		return fmt.Errorf("failed to delete zip file: %v", err)
	}
	Log.Infof("%s unzipped successfully and zip file deleted", zipFile)
	return nil
}
