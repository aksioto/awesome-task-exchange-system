package helper

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"log"
)

func EncryptAES(key []byte, plaintext string) string {
	c, err := aes.NewCipher(key)
	CheckError(err)

	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println("DECRYPTED:", s)
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}