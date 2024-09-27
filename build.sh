#!/bin/bash

#
# This script is used to build the project.
#
# The script takes the variables in the config.yml file and updates the project before building it.

# Parse the config.yml file

windows_workspace_path=$(yq e '.windows_workspace_path' config.yml)
linux_workspace_path=$(yq e '.linux_workspace_path' config.yml)
tools_folder_linux=$(yq e '.tools_folder_linux' config.yml)
tools_folder_windows=$(yq e '.tools_folder_windows' config.yml)
exploits_folder_linux=$(yq e '.exploits_folder_linux' config.yml)
exploits_folder_windows=$(yq e '.exploits_folder_windows' config.yml)
linpeas_linux=$(yq e '.linpeas_linux' config.yml)
linpeas_windows=$(yq e '.linpeas_windows' config.yml)

# Build the binaries for the project (amd64 only necessary 2024.09)
env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags='-extldflags=-static' -ldflags "-X 'cmd.LinToolsPath=$tools_folder_linux' -X 'cmd.LinExploitsPath=$exploits_folder_linux' -X 'cmd.Linpeas_file=$linpeas_linux'" -o build/oscptools_amd64-linux main.go

# Windows build
env GOOS=windows GOARCH=amd64 go build -ldflags "-X 'cmd.WinToolsPath=$tools_folder_windows' -X 'cmd.WinExploitsPath=$exploits_folder_windows' -X 'cmd.Winpeas_file=$linpeas_windows'" -o build/oscptools_amd64-windows.exe main.go

