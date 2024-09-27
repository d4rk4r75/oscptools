# OSCPTOOLS

## Usage
- Iniate by running linpeas/winpeas on targeted OS
- Quick download of multiple tools from locally hosted webserver (GoHTTPServer | python simplehttpserver, http.server)
- Upload files to specific share
- Mount SMB share on target. 


## Preconditions

#### HttpServer for downloading tools
```
python3 -m http.server <Port>
```

or 

```
~/.local/share/go/bin/gohttpserver -r /mnt/hgfs/VMToolsPentest --port 80 --upload
```

##### Directories

- /Tools
  - /windows
  - /linux
- /Exploits
  - /windows
  - /linux

The tool downloads all files recursively from the above directories, and stores them in "tools" and "exploits" in:
[Windows]
**Path**: C:\workspace\

[Linux]
**Path**: /tmp/workspace [all tools will have the execute bit set]


#### SMB Share with Impacket-SmbServer, preferrably in a folder in your loot directory
```
impacket-smbserver -smb2support -username <user> -password <pass> Public SMBUpload
```

#### HttpServer for uploading files in your loot directory
```
python3 -m uploadserver 8888
```

or 

```
~/.local/share/go/bin/gohttpserver -r /mnt/hgfs/VMToolsPentest --port 80 --upload
```
