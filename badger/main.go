package main

import (
	"fmt"
	"strconv"

	"github.com/sanity-io/litter"
	"github.com/satori/go.uuid"
	"github.com/tsingson/gin/fasthttputils"
	"github.com/tsingson/vk/txn"
)

var (
	tsn *txn.Txn
	err error
)

func main() {

	path, _ := fasthttputils.GetCurrentExecDir()

	tsn, err = txn.NewTxn(path)
	if err != nil {

	}
	defer tsn.Close()

	key := "string"

	for i := 0; i < 10000; i++ {

		key := []byte(key + strconv.Itoa(i))
		// set

		vvv := []byte(uuid.NewV4().String())
		err = tsn.Set(key, vvv)
		if err != nil {
			panic(err)
		}

		//	tsn.Db.Close()
		/**
		  err = tsn.Delete(key)

		  if err != nil {
		  	panic(err)
		  }
		*/
	}

	for i := 0; i < 10000; i++ {

		key := []byte(key + strconv.Itoa(i))

		// get

		value, err1 := tsn.GetStr(key)
		if err1 != nil {
			panic(err)
		}
		fmt.Println("")
		fmt.Println("")
		litter.Dump(value)

		//	tsn.Db.Close()
		/**
		  err = tsn.Delete(key)

		  if err != nil {
		  	panic(err)
		  }
		*/
	}
}
