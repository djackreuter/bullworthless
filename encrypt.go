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

	test := []byte("Test Secret stuff")
	fmt.Println(string(test))
	key := make([]byte, 32)
	fmt.Println(key)

	home, _ := os.UserHomeDir()
	fmt.Println("home ", home)

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	traverseFiles(files, ".")

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

	data := gcm.Seal(nonce, nonce, test, nil)
	fmt.Println(data)
	os.WriteFile("data.enc", data, 0777)
}

func dirRecurse(dirname os.DirEntry, base string) {
	fmt.Println("dir name ", dirname.Name())
	fmt.Println("base and name ", base+dirname.Name())
	files, err := os.ReadDir(base)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("there")
	traverseFiles(files, base)
}

func traverseFiles(files []os.DirEntry, base string) {
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name(), " is a dir")
			base = base + "\\" + file.Name()
			fmt.Println("base ", base)
			dirRecurse(file, base)
		} else {
			base = base + "\\"
			fmt.Println("just a file ", base+file.Name())
			data, err := os.ReadFile(base + file.Name())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("file data: ", data)
		}
	}
}
