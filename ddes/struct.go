package main

/*
#include <stdlib.h>
#include <string.h>

struct MyString
{
    char* s;
    int len;
};

struct MyString xmalloc(int len)
{
    static const char* s = "0123456789";
    char* p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    struct MyString str;
    str.s = p;
    str.len = len;
    return str;
}
*/
import "C"
import "unsafe"
import "fmt"

func main() {
	len := 10
	str := C.xmalloc(C.int(len))
	defer C.free(unsafe.Pointer(str.s))
	gostr := C.GoStringN(str.s, str.len)
	fmt.Printf("retlen=%v\n", str.len)
	println(gostr)
}
