// HIDEME: a tool for hiding large amounts of data on your computer using AES encryption
//
// notice: if someone does something bad with this its not my fault
// and also the hiding functionality is for really stupid people so someone might find out :|


package main

import (
	"math/rand"
	"math"
	"runtime"
	"os"
	"strings"
	"errors"
	"fmt"
	"time"
)


// probably a good idea to wipe this file afterwards
type hidemeConfig struct {
	randomFilenamesWin   []string `json:"randomFilenamesWin"`
	randomFilenamesUnix  []string `json:"randomFilenamesUnix"`
	randomExtensionsWin  []string `json:"randomExtensionsWin"`
	randomExtensionsUnix []string `json:"randomExtensionsUnix"`
	locationsTryWin      []string `json:"locationsTryWin"`
	locationsTryUnix     []string `json:"locationsTryUnix"`
	chunks               int      `json:"chunks"`
	passphrase           string   `json:"passphrase"`
}


// voodoo numbers
const (
	MagicDB = "SQLite format 3\000"
	MagicPE = "PK"
)


// get a random file name
func getFileName(config *hidemeConfig) string {
	if runtime.GOOS == "windows" {
                // window more like glass
		extension := "." + config.randomExtensionsWin[rand.Intn(len(config.randomExtensionsWin))]
                return config.randomFilenamesWin[rand.Intn(len(config.randomFilenamesWin))] + extension
        }

	// literally everything else
	extension := "." + config.randomExtensionsUnix[rand.Intn(len(config.randomExtensionsUnix))]
	return strings.ToLower(config.randomFilenamesUnix[rand.Intn(len(config.randomFilenamesUnix))]) + extension
}


// load the config file
func loadConfigFile(configPath string) *hidemeConfig {
	// default configuration
	config := &hidemeConfig{
		randomFilenamesWin:   []string{"authdb", "tempDiagnostic", "Info", "SysInfo", "SysDiagnostic", "SysAuth", "meta", "metaInfo", "tempfile", "taskProcessMeta", "FileData", // dumb system files
			"WData", "FData", "ITemp", "tempLog", "manifest", "infoLog", "metaData", "SysLog", "authConfig", "settings", "dump", "dumpTrace", "stackdump", "stacktrace", // more dumb system files
		},
		randomFilenamesUnix:  []string{"authdb", "tempDiagnostic", "SysDiagnostic", "Sysinfo", "meta", "metaInfo", "tempfile", "taskProcessMeta", "FileData", // dumb system files
                        "WData", "FData", "ITemp", "tempLog", "manifest", "infoLog", "metaData", "SysLog", "authConfig", "settings", "dump", "dumpTrace", "stackdump", "stacktrace", // more dumb system files
                },
		randomExtensionsWin:  []string{"db", "tmp", "sys", "plist", "log", "info", "inf", "dll", "pkg", "spec", "mta", "dat", "bin"},
		randomExtensionsUnix: []string{"db", "tmp", "log", "info", "inf", "pkg", "spec", "mta", "dat", "bin"},
		locationsTryWin:      []string{"c:/", "c:/temp", "c:/tmp", "c:/windows/Temp", "c:/users", "c:/users/%username%/AppData/", "c:/users/%username%/appdata/local",
			"c:/logs", "c:/perflogs", "c:/%WinLogs", "c:/riot games"}, // oh yes i sure do love people putting things in my valorant install folder
		locationsTryUnix:     []string{"/home/%username%/.local", "/home/%username%/.sysproc", "/home/%username%/", "/etc/.configdata/", "/etc/.tempinfo/"}, // hidden linux stuff
		chunks:               0, // auto-calculate the value like a big boy
	}

	return config
}

// hideme
func hideme(config *hidemeConfig, file string) error {
	// the main hideme function

	// seed the magic man
	rand.Seed(time.Now().UnixNano())

	// ensure the file is valid
	fileInfo, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("file doesn't exist")
	}

	if fileInfo.Size() == 0 {
		return errors.New("too small")
	}

	// see if we should calculate the split value
	if config.chunks == 0 {
		// the average size should be about 100 megabytes
		if fileInfo.Size() > 104857600 {
			config.chunks = int(math.Ceil(float64(fileInfo.Size()) / float64(104857600)))
		}
		config.chunks = 1
	}

	// get the filenames for each chunk
	var filenames []string
	for i := 0; i < config.chunks; i++ {
		filenames = append(filenames, getFileName(config))
	}

	fmt.Println(filenames)

	return nil
}

func main() {
	// main
	err := hideme(loadConfigFile("config.json"), "test")
	if err != nil {
		fmt.Printf(err.Error() + "\n")
	}
}
