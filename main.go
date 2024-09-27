// Description: Main file for oscptools.
//

package main

import (
	"flag"
	"os"

	"github.com/d4rk4r75/oscptools/cmd"
	"github.com/sirupsen/logrus"
)

var verbose bool

func main() {
	flag.StringVar(&cmd.KaliMachine, "k", "192.168.45.215", "IP Address of the Attacker machine")
	flag.IntVar(&cmd.HttpServerPort, "P", 80, "Port of the HTTP server")
	flag.StringVar(&cmd.GoHTTPServerPath, "PP", "Workspace", "Override the default Upload Path on the Web Server")
	flag.StringVar(&cmd.SmbShare, "s", "Public", "Name of the SMB share, commonly: Public")
	flag.StringVar(&cmd.SmbShareUser, "u", "pj", "Username for the SMB share")
	flag.StringVar(&cmd.SmbSharePassword, "p", "pj", "Password for the SMB share")
	flag.StringVar(&cmd.SmbShareDriveWin, "d", "F:", "Override the default drive letter for the SMB share")
	flag.IntVar(&cmd.UploadServerPort, "UP", 0, "Override the default Upload Server Port")
	flag.StringVar(&cmd.FileUpload, "U", "", "Upload a file to the SMB share, path/to/the/file")
	flag.BoolVar(&verbose, "debug", false, "Enable debug mode")
	flag.Parse()

	if verbose == true {
		cmd.Log.SetLevel(logrus.DebugLevel)
	}

	cmd.GetOSDetails()
	cmd.CheckShare()
	cmd.PrintBanner()

	if cmd.FileUpload != "" {
		cmd.UploadFile(cmd.FileUpload, "")
		os.Exit(0)
	}

	cmd.AskForInput()
}
