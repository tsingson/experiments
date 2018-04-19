package main

import (
	"github.com/sanity-io/litter"
	"github.com/tsingson/vk/tnx"
)

func main() {

	tsn, err := tnx.NewTxn()
	if err != nil {

	}
	defer tsn.Close()

	key := []byte("key")
	// set
	/**
	vvv := []byte(uuid.NewV4().String())
	err = tsn.Set(key, vvv)
	if err != nil {
		panic(err)
	}
	*/
	// get

	value, err1 := tsn.GetStr(key)
	if err1 != nil {
		panic(err)
	}
	litter.Dump(value)
	/**
	err = tsn.Delete(key)

	if err != nil {
		panic(err)
	}
	*/
}
