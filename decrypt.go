package main

import (
	"crypto/aes"
	"crypto/cipher"
	"runtime"
	"fmt"
	"os"
)

var fileSep string

func main() {
	home, _ := os.UserHomeDir()
	fmt.Println("home ", home)

	//testdir := ".\\test"
	testdir := "./test"

	opSystem := runtime.GOOS
	fileSep = "/"
	if opSystem == "windows" {
		fileSep = "\\"
	}

	traverseFiles(testdir)

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
			fmt.Println("dir path: ", base)
		} else {
			base := path + fileSep + file.Name()
			
			fmt.Println("file to dec: ", base)
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

	key := []byte{}

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
