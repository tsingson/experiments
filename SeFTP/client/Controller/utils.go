package Controller

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

//GCMEncrypter is a function to encrypt data using AES-GCM
func GCMEncrypter(data []byte, key [32]byte, nonce []byte) []byte {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Println(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	return ciphertext
}

//GCMDecrypter is a function to decipher data using AES-GCM
func GCMDecrypter(encData []byte, key [32]byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		log.Println(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, encData, nil)

	return plaintext, err
}
