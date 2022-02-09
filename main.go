// HIDEME: a tool for hiding large amounts of data on your computer using AES encryption
//
// notice: if someone does something bad with this its not my fault
// and also the hiding functionality is for really stupid people so someone might find out :|


package main

import (
	"math/rand",
	"math",
	"runtime",
	"strings",
	"errors",
)


// probably a good idea to wipe this file afterwards
type hidemeConfig struct {
	randomFilenames  []string `json:"randomFilenames"`
	randomExtentions []string `json:"randomExtensions"`
	locationsTryWin  []string `json:"locationsTryWin"`
	locationsTryUnix []string `json:"locationsTryUnix"`
	chunks           int      `json:"chunks"`
	passphrase       string   `json:"passphrase"`
}

// since i don't want to mess with cross-platform builds
*defaultConfig hidemeConfig = &hidemeConfig{
	randomFilenames:  ["authdb", "tempDiagnostic", "Info", "SysInfo", "SysDiagnostic", "SysAuth", "meta", "metaInfo", "tempfile", "taskProcessMeta", "FileData", // dumb system files
		"WData", "FData", "ITemp", "tempLog", "manifest", "infoLog", "metaData", "SysLog", "authConfig", "settings", "dump", "dumpTrace" // more dumb system files
			  ],
	randomExtensions: ["db", "tmp", "sys", "plist", "log", "info", "inf", "dll", "pkg", "spec", "mta", "dat", "bin"],
	locationsTryWin:  ["c:/", "c:/temp", "c:/tmp", "c:/windows/Temp", "c:/users", "c:/users/%username%/AppData/", "c:/users/%username%/appdata/local",
		"c:/logs", "c:/perflogs", "c:/%WinLogs", "c:/riot games"], // oh yes i sure do love people putting things in my valorant install folder
	locationsTryUnix: ["/home/%username%/.local", "/home/%username%/.sysproc", "/home/%username%/", "/etc/.configdata/", "/etc/.tempinfo/"], // hidden linux stuff
	chunks:           0, // auto-calculate the value like a big boy
}

// voodoo numbers
const (
	MagicDB = "SQLite format 3\000"
	MagicPE = "PK"
)


// get a random file name
func getFileName(config *hidemeConfig) string {
	extension := "." + config.randomExtensions[rand.Intn(len(config.randomExtensions))]

	if runtime.GOOS == "linux" {
		// objectively, the best os
		return strings.ToLower(config.randomFilenames[rand.Intn(len(config.randomFilenames))]) + extension
	}

	if runtime.GOOS == "windows" {
		// window more like glass
		return config.randomFilenames[rand.Intn(len(config.randomFilenames))] + extension
	}

	if runtime.GOOS == "darwin" {
		// yucky
		return config.randomFilenames[rand.Intn(len(config.randomFilenames))] + extension
	}
}

// load the config file
func loadConfigFile(configPath string) *hidemeConfig {
	config := defaultConfig

	
}

// hideme
func hideme(config *hidemeConfig, file string) error {
	// the main hideme function

	// ensure the file is valid
	fileInfo, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("file doesn't exist")
	}

	// see if we should calculate the split value
	if config.chunks == 0 {
		// the average size should be about 100 megabytes
		if fileInfo.Size() > 104857600 {
			config.chunks = math.Ceil(fileInfo.Size() / 104857600)
		}
	}

	// get the filenames for each chunk
	var filenames []string
	for i := 0; i < config.chunks; i++ {
		filenames = append(filenames, getFileName(config))
	}
}
