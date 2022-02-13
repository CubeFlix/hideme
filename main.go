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
	"path"
	"strings"
	"errors"
	"fmt"
	"bufio"
	"time"
	"os/user"
	"golang.org/x/crypto/pbkdf2"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha1"
	"io"
	"encoding/binary"
)


// dumb version info
const (
	Version = "v0.0"
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
	chunkSize            int64    `json:"chunkSize"`
	password             string   `json:"password"`
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
	return "." + strings.ToLower(config.randomFilenamesUnix[rand.Intn(len(config.randomFilenamesUnix))]) + extension
}


// get a random path
func getHidePath(config *hidemeConfig, username string) string {
	hidepath := ""

	if runtime.GOOS == "windows" {
                // window more like glass
		hidepath = config.locationsTryWin[rand.Intn(len(config.locationsTryWin))]
        } else {
		// literally everything else
		hidepath = config.locationsTryUnix[rand.Intn(len(config.locationsTryUnix))]
	}

	hidepath = strings.Replace(hidepath, "%username%", username, -1)
	return hidepath
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
			"c:/logs", "c:/perflogs", "c:/%WinLogs", "c:/Riot Games"}, // oh yes i sure do love people putting things in my valorant install folder
		locationsTryUnix:     []string{"/home/%username%/.local", "/home/%username%/.sysproc", "/home/%username%/", "/etc/.configdata/", "/etc/.tempinfo/"}, // hidden linux stuff
		chunks:               0, // auto-calculate the value like a big boy
		chunkSize:            52428800,
	}

	return config
}

// hide the files
func hidefiles(config *hidemeConfig, paths []string, file string) error {
	fmt.Println()
	fmt.Println("beginning encryption process")

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	// create the AES encryption
	dk := pbkdf2.Key([]byte(config.password), salt, 4096, 32, sha1.New)
	fmt.Println("generated key")
	block, err := aes.NewCipher(dk)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return err
	}
	fmt.Println(nonce, dk, gcm, block)
	fmt.Println("generated block")

	// open the file
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)

	// sepreately encrypt and write each chunk
	// go through each chunk file
	for i := 0; i < len(paths); i++ {
		fmt.Println("generating chunk:", i, paths[i])
		// read the section
		extraBytes := 4
		if i == len(paths) - 1 {
			// last one
		} else {
			// account for the length of the next path
			extraBytes += len(paths[i + 1])
		}
		buf := make([]byte, 4)
		// include extra bytes
		if i != len(paths) - 1 {
			binary.LittleEndian.PutUint32(buf, uint32(len(paths[i + 1])))
			buf = append(buf, []byte(paths[i + 1])...)
		} else {
			binary.LittleEndian.PutUint32(buf, uint32(0))
		}

		// start actual data
		newBuf := make([]byte, config.chunkSize)
		n, err := reader.Read(newBuf) // loading chunk into buffer
		if err != nil {
			return err
		}
		newBuf = newBuf[:n]
		buf = append(buf, newBuf...)

		// create the new file
		chunkFile, err := os.Create(paths[i])
		if err != nil {
			return err
		}

		// write m a g i c
		if strings.HasSuffix(paths[i], ".db") {
			chunkFile.Write([]byte(MagicDB))
		} else if strings.HasSuffix(paths[i], ".exe") || strings.HasSuffix(paths[i], ".sys") || strings.HasSuffix(paths[i], ".dll") {
			chunkFile.Write([]byte(MagicPE))
		}

		chunkFile.Write(salt)

		dstNonce := make([]byte, gcm.NonceSize())
		dstNonce = nonce
		// we can encrypt the buffer
		encryptedData := gcm.Seal(dstNonce, nonce, buf, nil)
		bytesWritten, err := chunkFile.Write(encryptedData)
		if err != nil {
			chunkFile.Close()
			return err
		}
		fmt.Println(len(encryptedData))
		fmt.Println(encryptedData[:64], encryptedData[len(encryptedData) - 64:])
		fmt.Println("wrote", bytesWritten, "bytes")
		chunkFile.Close()
	}

	fmt.Println("done")
	return nil
}


// decrypt the data
func decryptfiles(path string, password string, save string) error {
	fmt.Println()

	// get the path
	if path == "" {
		fmt.Printf("no path included. please enter the name of the main path: ")
		fmt.Scanln(&path)
	}

	// get the password
	if password == "" {
		fmt.Printf("password for decryption: ")
		fmt.Scanln(&password)
	}

	// get the save path
	if save == "" {
		fmt.Printf("save path: ")
		fmt.Scanln(&save)
	}

	// create the aes decryptor
        fmt.Println("beginning decryption process")

	// read out the salt from the first file
	firstFile, err := os.Open(path)
        if err != nil {
		return err
	}
	ffReader := bufio.NewReader(firstFile)
	// handle magic
	salt := make([]byte, 16)
	if strings.HasSuffix(path, ".db") {
		ffbuf := make([]byte, len(MagicDB) + 16)
                ffReader.Read(ffbuf)
		salt = ffbuf[len(ffbuf) - 16:]
        } else if strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".sys") || strings.HasSuffix(path, ".dll") {
		ffbuf := make([]byte, len(MagicPE) + 16)
		ffReader.Read(ffbuf)
		salt = ffbuf[len(ffbuf) - 16:]
	} else {
		ffbuf := make([]byte, 16)
		ffReader.Read(ffbuf)
		salt = ffbuf[len(ffbuf) - 16:]
	}

        // create the AES encryption
        dk := pbkdf2.Key([]byte(password), salt, 4096, 32, sha1.New)
        fmt.Println("generated key")
        block, err := aes.NewCipher(dk)
        if err != nil {
		firstFile.Close()
                return err
        }
        gcm, err := cipher.NewGCM(block)
        if err != nil {
		firstFile.Close()
                return err
        }
        nonce := make([]byte, gcm.NonceSize())
        ffReader.Read(nonce)
	firstFile.Close()
	fmt.Println(nonce, dk, gcm, block)
	fmt.Println("generated block")

        // start the loop (for each chunk file)
	currentpath := path
	for {
		fileinfo, err := os.Stat(currentpath)
		if err != nil {
			return err
		}
		fileSize := fileinfo.Size()
		reader, err := os.Open(currentpath)
		if err != nil {
			return err
		}

		// remove beginning bytes
		// handle magic
		extraBytes := 0
		if strings.HasSuffix(currentpath, ".db") {
			extraBytes = len(MagicDB) + 16 + gcm.NonceSize()
			fileSize -= int64(len(MagicDB)) + 16 + int64(gcm.NonceSize())
		} else if strings.HasSuffix(currentpath, ".exe") || strings.HasSuffix(currentpath, ".sys") || strings.HasSuffix(currentpath, ".dll") {
			extraBytes = len(MagicPE) + 16 + gcm.NonceSize()
			fileSize -= int64(len(MagicPE)) + 16 + int64(gcm.NonceSize())
		} else {
			fmt.Println(16 + gcm.NonceSize())
			extraBytes = 16 + gcm.NonceSize()
			fileSize -= 16 + int64(gcm.NonceSize())
		}

		fmt.Println(salt, nonce)

		// decrypt
		encData := make([]byte, fileSize + int64(extraBytes))
		fmt.Println(fileSize)
		_, err = reader.Read(encData)
		if err != nil {
			reader.Close()
			return err
		}
		reader.Close()
		encData = encData[extraBytes:]
		fmt.Println(encData[:64], encData[len(encData) - 64 : ])
		data, err := gcm.Open(nonce, nonce, encData, nil)
		if err != nil {
			reader.Close()
			return err
		}

		fmt.Println(data[0 : 16])
		fmt.Println(string(data[:32]))
		break
	}
	return nil
}

// hideme
func hideme(config *hidemeConfig, file string) error {
	// the main hideme function

	// h e a d e r
	fmt.Println("==========")
	fmt.Println("hideme " + Version)
	fmt.Println()
	fmt.Println("hide all your files in plain sight")
	fmt.Println("i am not responsible for illegal activity caused by this tool")
	fmt.Println("==========")
	fmt.Println()

	// seed the magic man
	rand.Seed(time.Now().UnixNano())

	// ensure the file is valid
	fileInfo, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("file doesn't exist")
	}

	if fileInfo.Size() == 0 {
		return errors.New("file too small")
	}

	// see if we should calculate the split value
	fmt.Println("file size:", fileInfo.Size(), "chunks:", math.Ceil(float64(fileInfo.Size()) / float64(config.chunkSize)))
	if config.chunks == 0 {
		// the average size should be about 100 megabytes
		if fileInfo.Size() > config.chunkSize {
			config.chunks = int(math.Ceil(float64(fileInfo.Size()) / float64(config.chunkSize)))
		} else {
			config.chunks = 1
		}
	}

	// get the filenames and paths for each chunk
	var paths []string
	for i := 0; i < config.chunks; i++ {
		// get a path to store in
		// resolve username
		userobj, err := user.Current()
		username := ""
		if err != nil {
			fmt.Printf("username for resolving user path: ")
			fmt.Scanln(&username)
		} else {
			username = userobj.Username
		}
		storepath := getHidePath(config, username)

		for {
			if os.MkdirAll(storepath, os.ModePerm) == nil {
				break
			}
			fmt.Println(os.MkdirAll(storepath, os.ModePerm).Error())
			storepath = getHidePath(config, username)
		}

		paths = append(paths, path.Join(storepath, getFileName(config)))
	}

	fmt.Println("hiding in: ", paths)

	// get a password if required
	if config.password == "" {
		fmt.Printf("password for encryption: ")
		fmt.Scanln(&config.password)
	}

	// create the new files
	err = hidefiles(config, paths, file)

	return err
}


// decryption
func decryptme(path string, password string, file string) error {
        // the main hideme function

        // h e a d e r
        fmt.Println("==========")
        fmt.Println("hideme " + Version)
        fmt.Println()
        fmt.Println("hide all your files in plain sight")
        fmt.Println("i am not responsible for illegal activity caused by this tool")
        fmt.Println("==========")
        fmt.Println()

        // seed the magic man
        rand.Seed(time.Now().UnixNano())

	// do the work
	err := decryptfiles(path, password, file)
	return err
}


func main() {
	// main
	// err := hideme(loadConfigFile("test"), "200MB.bin")
	err := decryptme("", "", "")
	if err != nil {
		fmt.Printf(err.Error() + "\n")
	}
}
