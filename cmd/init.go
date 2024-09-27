package cmd

import (
	"fmt"
)

func PrintBanner() {
	// Description: Print the banner.
	fmt.Println("_____________________________________________________________________")
	fmt.Println("╔═╗╔═╗╔═╗╔═╗  ╔╦╗╔═╗╔═╗╦  ╔═╗")
	fmt.Println("║ ║╚═╗║  ╠═╝   ║ ║ ║║ ║║  ╚═╗")
	fmt.Println("╚═╝╚═╝╚═╝╩     ╩ ╚═╝╚═╝╩═╝╚═╝")
	fmt.Println("_____________________________________________________________________\n")

	fmt.Printf("Attacker Tools:\t http://%s:%d/ [%s]\n", KaliMachine, HttpServerPort, CheckHttpServerPort())
	fmt.Printf("SMB Share:\t %s [%s]\n", SmbShareFull, SmbShareStatus)
	fmt.Printf("Upload Server:\t %s [%s]\n", UploadServerFull, UploadServerStatus)

	fmt.Println("_____________________________________________________________________\n")
	fmt.Printf("OS:\t\t %s\n", OS)
	fmt.Printf("User:\t\t %s %s\n", CurrentUser, CurrentGroups)
	fmt.Println("_____________________________________________________________________\n")
}
