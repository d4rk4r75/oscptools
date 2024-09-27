package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string, lpath string, fileName string) (string, error) {
	// Description: Download a file from the internet
	fPath := filepath.Join(lpath, fileName)

	Log.Debugf("Downloading file from: %s\n", url)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		Log.Errorf("Error making GET request: %v\n", err)
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fPath)
	if err != nil {
		Log.Errorf("Error creating file: %v\n", err)
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		Log.Errorf("Error writing to file: %v\n", err)
		return "", fmt.Errorf("error writing to file: %v", err)
	}

	return fPath, nil
}
