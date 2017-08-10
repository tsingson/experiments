package main

import (
	"encoding/base64"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	srcKey := "1234567890qwerty"

	srcMsg := "1234567890abcdefg-^&*KillBill----1234567890abcdefg-^&*KillBill----"
	aesEnc := NewAesEncrypt(srcKey)
	arrEncrypt, err := aesEnc.Encrypt(srcMsg)
	if err != nil {
		//	arrEncryptString := base64.StdEncoding.EncodeToString(arrEncrypt)
		fmt.Println(arrEncrypt)
		return
	}
	spew.Dump(arrEncrypt)
	arrEncryptString := base64.StdEncoding.EncodeToString(arrEncrypt)
	fmt.Println(arrEncryptString)

	// data, err := base64.StdEncoding.DecodeString(string(text))
	strMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		fmt.Println(arrEncrypt)
		return
	}
	fmt.Println(strMsg)
}
