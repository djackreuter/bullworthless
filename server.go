package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"crypto/sha256"
	"encoding/pem"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {

	var genRsaKeys bool
	var encAesKey string

	flag.BoolVar(&genRsaKeys, "genKeys", false, "Generate server RSA keys")
	flag.StringVar(&encAesKey, "decryptAES", "", "Decrypt targets AES key")
	flag.Parse()

	if genRsaKeys == true {
		err := genKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if encAesKey != "" {
		aesKey, err := decryptAesKey(encAesKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(aesKey)
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

func genKeys() error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pubKey := &privKey.PublicKey
	return writeKeys(privKey, pubKey)
}

func decryptAesKey(hexKey string) ([]byte, error) {
	var aesKey []byte
	aesEncKey, err := hex.DecodeString(hexKey)
	if err != nil {
		fmt.Println(err)
		return aesKey, err
	}
	privKeyFile, err := os.Open("private_key.pem")
	if err != nil {
		fmt.Println(err)
		return aesKey, err
	}
	pemfileinfo, _ := privKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privKeyFile.Close()

	privKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		return aesKey, err
	}
	aesKey, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, aesEncKey, nil)
	if err != nil {
		fmt.Println(err)
		return aesKey, err
	}
	fmt.Println(aesKey)
	return aesKey, nil
}

