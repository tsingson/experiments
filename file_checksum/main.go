package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"reflect"
	"unsafe"
)

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func shaCheckSum(filename string) (result string, err error) {
	f, f_err := os.Open(filename)
	if f_err != nil {
		log.Fatal(f_err)
		f.Close()
		return "", f_err
	}
	defer f.Close()

	h := sha256.New()
	if _, s_err := io.Copy(h, f); s_err != nil {
		log.Fatal(s_err)
	}

	return BytesToStringUnsafe(h.Sum(nil)), nil
}

// BytesToStringUnsafe converts a byte slice to a string.
// It's fast, but not safe. Use it only if you know what you're doing.
func BytesToStringUnsafe(b []byte) string {
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: bytesHeader.Data,
		Len:  bytesHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}

// StringToBytesUnsafe converts a string to a byte slice.
// It's fast, but not safe. Use it only if you know what you're doing.
func StringToBytesUnsafe(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bytesHeader := reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bytesHeader))
}

func main() {

}
