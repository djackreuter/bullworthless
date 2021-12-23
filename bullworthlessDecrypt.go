package main

import (
	"crypto/aes"
	"crypto/cipher"
	"runtime"
	"encoding/hex"
	"fmt"
	"flag"
	"os"
)

var fileSep string
var key []byte

func main() {
	home, _ := os.UserHomeDir()

	var encKey string
	var err error

	flag.StringVar(&encKey, "key", "", "Decrypted hex key")
	flag.Parse()

	opSystem := runtime.GOOS
	fileSep = "/"
	if opSystem == "windows" {
		fileSep = "\\"
	}
	if encKey == "" {
		fmt.Println("Need decrypted key")
		os.Exit(1)
	}
	key, err = hex.DecodeString(encKey)
	if err != nil {
		fmt.Println("Error decoding key: ", err)
	}
	traverseFiles(home)
}

func traverseFiles(path string) {
	var dirs []string

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			base := path + fileSep + file.Name()
			dirs = append(dirs, base)
		} else {
			base := path + fileSep + file.Name()
			decryptFile(base)
		}
	}
	for _, dir := range dirs {
		traverseFiles(dir)
	}
}

func decryptFile(path string) {

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		fmt.Println(err)
		os.Exit(1)
	}

	nonce := data[:nonceSize]
	encData := data[nonceSize:]
	decData, err := gcm.Open(nil, nonce, encData, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(path, decData, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

