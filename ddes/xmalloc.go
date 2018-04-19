package main

/*
#include <stdlib.h>
#include <string.h>
char* xmalloc(int len, int *rlen)
{
    static const char* s = "0123456789";
    char* p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    *rlen = len;
    return p;
}
*/
import "C"
import "unsafe"
import "fmt"

func main() {
	rlen := C.int(0)
	len := 10
	cstr := C.xmalloc(C.int(len), &rlen)
	defer C.free(unsafe.Pointer(cstr))
	gostr := C.GoStringN(cstr, rlen)
	fmt.Printf("retlen=%v\n", rlen)
	println(gostr)
}
