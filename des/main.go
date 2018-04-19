package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"reflect"
	"unsafe"

	"github.com/sanity-io/litter"
)

func main() {

	key := []byte{0xa3, 0xbe, 0x93, 0xff, 0x10, 0x34, 0x5f, 0xde,
		0xc6, 0x2e, 0x57, 0x83, 0x29, 0x7c, 0x8e, 0xf6,
		0xa3, 0x58, 0x34, 0x27, 0x13, 0x2c, 0x4e, 0xd2}

	iv := []byte("34234323")
	data := []byte(`{"activationcode":"123445678asd","mac":"00:50:56:31:40:38","mac_2":"03:40:76:A1:4C:E8","apktype":"WonderfulOttMovie","libversion":"libVK_STBProxy1.3.9","nativesn":"12.00-09.10-10000000","platforminfo":"Android 6.01","grade":3}`)
	out, _ := DesEncrypt(data, key, iv)
	debyte := base64Encode(out)
	litter.Dump(string(debyte))

	/**
	fmt.Println(debyte)

	enbyte, _ := base64Decode(debyte)
	litter.Dump(enbyte)

	//log.Println("base:", enbyte)
	out, _ = DesDecrypt(enbyte, key, iv)
	litter.Dump(out)
	*/
}

func BytesToStringUnsafe(b []byte) string {
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: bytesHeader.Data,
		Len:  bytesHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
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
