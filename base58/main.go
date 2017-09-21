package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/satori/go.uuid"
	"github.com/tsingson/fastweb/guid"
)

func main() {
	enc := new(guid.Base58)
	u := uuid.NewV4()
	fmt.Println("uuid version 4", u.String())
	uu := hex.EncodeToString(u.Bytes())

	fmt.Println("uuid version ----", uu)

	/**
	shortuid := enc.UuidEncode(uu)
	fmt.Println("uuid 4 base58", shortuid) // 6R7VqaQHbzC1xwA5UueGe6

	uid, _ := enc.UuidDecode(shortuid)
	fmt.Println(uid)
	*/
	ip := "192.168.1.1"
	fmt.Println("ip", ip)
	ips := enc.Encode([]byte(ip))

	fmt.Println("ip be base58 encode", ips)
	fmt.Println("ip be base58 decode", string(enc.Decode(ips)))

	fmt.Println("****************")
	res := base58.CheckEncode([]byte(ip), 20)
	fmt.Println("base58check", res)
	orgin, _, _ := base58.CheckDecode(res)
	fmt.Println("base58check decode ", string(orgin))

	fmt.Println("****************")
	uures := base58.Encode([]byte(uu))
	fmt.Println("base58", uures)
	uuorgin := base58.Decode(uures)
	fmt.Println("base58 decode ", string(uuorgin))

}
