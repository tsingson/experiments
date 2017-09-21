package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var mdata string = "hellow world , kill 1234567890"
	spew.Dump(mdata)
	key := []byte("76543456")
	iv := []byte("34234323")
	data := []byte(mdata)
	out, _ := DesEncrypt(data, key, iv)
	debyte := base64Encode(out)
	spew.Dump(debyte)
	fmt.Println(debyte)

	enbyte, _ := base64Decode(debyte)
	//log.Println("base:", enbyte)
	out, _ = DesDecrypt(enbyte, key, iv)
	spew.Dump(out)
}
func DesEncrypt(data, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	out := make([]byte, len(data))
	blockMode.CryptBlocks(out, data)
	return out, nil
}
func DesDecrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	out := make([]byte, len(data))
	blockMode.CryptBlocks(out, data)
	out = PKCS5UnPadding(out)
	return out, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
