package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)



func main() {

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
			decryptFile(data, path)
			fmt.Println("file data: ", data)
		}
	}
}


func decryptFile(data []byte, path string) {
	key := make([]byte, 32)

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

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("path pre chop ", path)
	path = path[:len(path)-4]
	fmt.Println("path post chop ", path)
	err = os.WriteFile(path, plaintext, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plaintext))
}
