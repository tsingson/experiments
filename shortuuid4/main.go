package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"

	"github.com/renstrom/shortuuid"
)

func FlateEncode(input string) (result []byte, err error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, -1)
	w.Write([]byte(input))
	w.Close()
	result = buf.Bytes()
	return
}

func FlateDecode(input []byte) (result []byte, err error) {
	result, err = ioutil.ReadAll(flate.NewReader(bytes.NewReader(input)))
	return
}

func main() {
	u := shortuuid.New() // Cekw67uyMpBGZLRP2HFVbe
	fmt.Println(u)
	fmt.Println(len(u))
	fmt.Println("Cekw67uyMpBGZLRP2HFVbe")
	fmt.Println(len("Cekw67uyMpBGZLRP2HFVbe"))
	fmt.Println(len("152193295848570880"))
}
