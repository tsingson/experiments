package main

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renstrom/shortuuid"
	"github.com/satori/go.uuid"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u.Bytes())
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}

func main() {
	enc := base58Encoder{}
	u := shortuuid.NewWithEncoder(enc)
	fmt.Println(u) // 6R7VqaQHbzC1xwA5UueGe6
	fmt.Println(len(u))
	fmt.Println("Ej5FZ90ibEtOkVkJmVUQAAA")
	fmt.Println(len("Ej5FZ90ibEtOkVkJmVUQAAA"))
	uu := uuid.NewV4()
	var k string = hex.EncodeToString(uu.Bytes())
	fmt.Println(k)
	fmt.Println(len(k))
	sb := uu.Bytes()
	fmt.Println("len*************", len(string(sb)))
	hexString := hex.EncodeToString(sb)
	hexByte, err := hex.DecodeString(hexString)
	fmt.Println(hexString)
	// 68656c6c6f20776f726c6421

	fmt.Println(hexByte, err)
	// [104 101 108 108 111 32 119 111 114 108 100 33] <nil>

	base32StdString := base32.StdEncoding.EncodeToString(sb)
	base32HexString := base32.HexEncoding.EncodeToString(sb)
	fmt.Println("uuid", uu.String())
	fmt.Println("base32", base32StdString)
	fmt.Println(len(base32StdString))
	// NBSWY3DPEB3W64TMMQQQ====

	fmt.Println("base32hex", base32HexString)
	// D1IMOR3F41RMUSJCCGGG====
	fmt.Println(len("0pPKHjWprnVxGH7dEsAoXX2YQvU"))
	fmt.Println(len("152193295848570880"))
	fmt.Println(len("123456789abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"))
}
