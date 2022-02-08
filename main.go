// HIDEME: a tool for hiding large amounts of data on your computer using AES encryption
//
// notice: if someone does something bad with this its not my fault
// and also the hiding functionality is for really stupid people so someone might find out :|


package main

// probably a good idea to wipe this file afterwards
type hidemeConfig struct {
	randomFilenames  []string `json:"randomFilenames"`
	randomExtentions []string `json:"randomExtensions"`
	locationsTryWin  []string `json:"locationsTryWin"`
	locationsTryUnix []string `json:"locationsTryUnix"`
	split            int      `json:"split"`
	passphrase       string   `json:"passphrase"`
}

defaultConfig hidemeConfig = hidemeConfig{
	randomFilenames:  ["Thumbs", "auth", "tempDiagnostic", "Info", "SysInfo", "SysDiagnostic", "SysAuth", "meta", "metaInfo", "tempfile", "taskProcessMeta", "FileData", // dumb system files
		"WData", "FData", "ITemp", "tempLog", "manifest", "infoLog", "metaData", "SysLog", "authConfig", "settings", "dump", "dumpTrace" // more dumb system files
			  ],
	randomExtensions: ["db", "tmp", "sys", "plist", "log", "info", "inf", "dll", "pkg", "spec", "mta", "dat", "bin"],
	locationsTryWin:  ["c:/", "c:/temp", "c:/tmp", "c:/windows/Temp", "c:/users", "c:/users/%username%/AppData/", "c:/users/%username%/appdata/local",
		"c:/logs", "c:/perflogs", "c:/%WinLogs", "c:/riot games"] // oh yes i sure do love people putting things in my valorant install folder
	locationsTryUnix: ["/home/%username%/.local", "/home/%username%/.sysproc", "/home/%username%/", "/etc/.configdata/", "/etc/.tempinfo/"] // hidden linux stuff
}


// get a random file name
func getFileNames

// load the config
func loadConfigFile
