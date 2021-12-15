package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
)



func main() {
	key := make([]byte, 32)
	//key := []byte("test")

	ciphertext, err := os.ReadFile("data.enc")
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
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
		os.Exit(1)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(plaintext))

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
			decryptFile(data, path)
			fmt.Println("file data: ", data)
		}
	}
}

func decryptFile(data []byte, path string) {

}
