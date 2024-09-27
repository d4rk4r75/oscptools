package cmd

var (
	// Description: These variables are used to store information about current system

	OS            string
	CurrentUser   string
	CurrentGroups []string
	IPAddress     string
	HostName      string

	// Tools Variables
	HttpServerFull   string
	HttpServerStatus string

	// SMB Share Variables
	SmbShareFull   string
	SmbShareStatus string

	UploadServerFull   string
	UploadServerPort   int
	UploadServerStatus string

	FileUpload string

	// OSCP Tools Variables
	KaliMachine      string
	HttpServerPort   int
	GoHTTPServerPath string
	SmbShare         string
	SmbShareUser     string
	SmbSharePassword string
	SmbShareDriveWin string
	SmbShareDriveLin string
)

// Tools Static Variables
const (
	WinToolsPath      = "/Tools/windows"
	WinExploitsPath   = "/Exploits/windows"
	LinuxToolsPath    = "/Tools/linux"
	LinuxExploitsPath = "/Exploits/linux"
	GoHTTPDownloadZip = "?op=archive"
	Linpeas_file      = "linpeas.sh"
	Winpeas_file      = "/PEASS-ng/winPEASx64.exe"
)

// Windows Static Variables
const (
	// Description: These variables are used to store static information about Windows systems
	WinTempPath   = "C:\\Workspace"
	LinuxTempPath = "/tmp/workspace/"
)
