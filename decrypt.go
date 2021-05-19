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
