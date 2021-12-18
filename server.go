package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pubKey := &privKey.PublicKey
	fmt.Println("pub key: ", pubKey)
	fmt.Println("priv key: ", privKey)

	// write to file
	pemPrivFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// start export process
	var pemPrivBlock = &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	// write block to file
	err = pem.Encode(pemPrivFile, pemPrivBlock)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPrivFile.Close()
}

func importKey() {
	privKeyFile, err := os.Open("private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemfileinfo, _ := privKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privKeyFile)
	_, err := buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privKeyFile.Close()

	privKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Private Key: ", privKeyImported)
	pubKey := &privKeyImported.PublicKey
	fmt.Println("Public Key: ", pubKey)
	fmt.Println("Private Key: ", privKeyImported)
}

