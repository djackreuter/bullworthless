# Bullworthless
 ## Ransomware impementation in Go
I wrote this to better understand how ransomeware and encryption works. Public and Private RSA keys are included in repo for easy POC and to discourage abuse. The attacker server generates a public and private RSA key with `server.go`, and the public key is then hardcoded into `bullworthless.go`.  File encryption is done via AES with a randomly generated 32 bit key. The AES key is then encrypted with the public key of the attacker server and then sent off to the attacker server.
> DISCLAIMER: This program will encrypt files on the computer on which it's ran!! It is imperative that end user follow all instructions to ensure safe recovery of files after encryption. Usage of this software without prior consent is illegal. It's the end user's responsibility to obey all applicable local, state, and federal laws. This software is provided "AS IS" and the Developer assumes no liability, and is not responsible for any misuse, data loss, or other damages that may be caused by this program. This software is for educational purposes only!

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
	* **Windows**:  `export GOOS=windows GOARCH=amd64;go build -ldflags "-w -s" bullworthless.go`
	* **Linux**:  `export GOOS=linux GOARCH=amd64;go build -ldflags "-w -s" bullworthless.go`
5. Run

### Decryption Steps
The attacker server or webhook will receive the encrypted AES key. Run: `./server -decryptAES <encrypted_key>` in the same directory as the `keys` folder.
This will get you back the AES key.

Then to decrypt the files again run: `./BullworthlessDecrypt -key <decrypted_key>` on the target machine.

---
You can find precompiled `server` and `bullworthlessDecrypt` binaries in the releases directory. `bullworthless.go` needs modification before compilation, so you need to compile that one yourself after making the aforementioned changes.

