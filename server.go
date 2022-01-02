package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"path/filepath"
	"crypto/sha256"
	"encoding/pem"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {

	var genRsaKeys bool
	var encAesKey string

	flag.BoolVar(&genRsaKeys, "genkeys", false, "Generate server RSA keys")
	flag.StringVar(&encAesKey, "decryptAES", "", "Decrypt targets AES key")
	flag.Parse()

	if genRsaKeys == true {
		err := genKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created public_key.pem and private_key.pem in keys dir")
	}

	if encAesKey != "" {
		aesKey, err := decryptAesKey(encAesKey)
		if err != nil {
			fmt.Println(err, ": Possibly missing or incorrect private key")
			os.Exit(1)
		}
		base64key := base64.RawStdEncoding.EncodeToString(aesKey)
		fmt.Println("Decrypted key: ", base64key)
	}
}

func writeKeys(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) error {
	keyDir := filepath.Join(".", "keys")
	err := os.MkdirAll(keyDir, 0755)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p := filepath.FromSlash("keys/private_key.pem")
	pemPrivFile, err := os.Create(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p = filepath.FromSlash("keys/public_key.pem")
	pemPubFile, err := os.Create(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var pemPrivBlock = &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	var pemPubBlock = &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}
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

func decryptAesKey(base64key string) ([]byte, error) {
	var aesKey []byte
	aesEncKey, err := base64.RawStdEncoding.DecodeString(base64key)
	if err != nil {
		fmt.Println(err)
		return aesKey, err
	}
	p := filepath.FromSlash("keys/private_key.pem")
	privKeyFile, err := os.Open(p)
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
		return aesKey, err
	}
	return aesKey, nil
}

