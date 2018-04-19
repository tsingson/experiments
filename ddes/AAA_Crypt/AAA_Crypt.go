package AAA_Crypt

/*
#include <stdint.h>
#include "AAA_Crypt.h"
#include "avdes.h"
#include "Base64.h"
#include "Des.h"
*/
import "C"

func GoDesEnCrypt(in_data_string string) *C.uint8_t {
	cs := C.CString(in_data_string)
	// defer C.Free(unsafe.Pointer(cs))
	return C.DesEnCrypt(cs)

}

func EnCrypt(in_data_string string) string {
	cs := C.CString(in_data_string)
	str := C.MyEnCrypt(cs)
	gostr := C.GoStringN(str.s, str.len)
	fmt.Printf("retlen=%v\n", str.len)
	return gostr
}

func GoTest() C.int {
	//cs := C.CString(in_data_string)
	// defer C.Free(unsafe.Pointer(cs))
	return C.Test()

}

// design and code by tsingson
