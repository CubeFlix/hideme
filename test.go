package main

import (
        "os"
        "fmt"
        "bufio"
	"crypto/cipher"
	"crypto/aes"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
)


const path = "/home/cubeflix/.local/.tempdiagnostic.mta"

func main() {
	fileinfo, err := os.Stat(path)
        if err != nil {
                fmt.Println(err.Error())
        }
        fileSize := fileinfo.Size()
	// fileSize -= 32
	f, err := os.Open(path)
	r := bufio.NewReader(f)
	//r.Discard(32)
	buf := make([]byte, fileSize)
	n, err := r.Read(buf)
	fmt.Println(n)
	f.Close()

	f, err = os.Open(path)
	n, err = f.Read(buf)
	fmt.Println(n)
	f.Close()

	salt := buf[:16]

	dk := pbkdf2.Key([]byte("hello"), salt, 4096, 32, sha1.New)
        fmt.Println("generated key")
        block, err := aes.NewCipher(dk)
        if err != nil {
		fmt.Println(err.Error())
                return
        }
        gcm, err := cipher.NewGCM(block)
        if err != nil {
                fmt.Println(err.Error())
                return
        }
	nonce := buf[16 : 16 + gcm.NonceSize()]
	fmt.Println(nonce, dk, gcm, block)
	data, err := gcm.Open(nil, nonce, buf[16 + gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)
}
