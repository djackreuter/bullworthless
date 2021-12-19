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

	privKey, pubKey := genKeys()

	err := writeKeys(privKey, pubKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeKeys(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) error {
	pemPrivFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println(err)
		return err
	}
	pemPubFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// start export process
	var pemPrivBlock = &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	var pemPubBlock = &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}
	// write block to file
	err = pem.Encode(pemPrivFile, pemPrivBlock)
	if err != nil {
		return err
	}
	pemPrivFile.Close()
	err = pem.Encode(pemPubFile, pemPubBlock)
	if err != nil {
		return err
	}
	pemPubFile.Close()
	return nil
}

func genKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pubKey := &privKey.PublicKey
	return privKey, pubKey
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
	_, err = buffer.Read(pembytes)

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

