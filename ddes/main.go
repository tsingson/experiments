package main

/**
#include <string.h>
#includ <stdlib.h>

char * GetChar()

char* GetChar() {

}
*/

import (
	"github.com/sanity-io/litter"
	"github.com/tsingson/experiments/ddes/AAA_Crypt"
)

func main() {

	url := `1234567890asdfghj`
	//AAA_Crypt.AaaEncrypt(url, len(url), out_data, 1024)
	litter.Dump(url)
	ok := AAA_Crypt.GoDesEnCrypt(url)
	litter.Dump(ok)
	i := AAA_Crypt.GoTest()
	litter.Dump(i)
}

// design and code by tsingson
