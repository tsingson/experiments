package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>


struct t {
char *s;
};

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var t C.struct_t
	var s string = `hello world 中国人民`
	var ch *C.char
	var tmp *C.char

	// 分配空间, 并判断是否分配成功
	t.s = (*C.char)(C.malloc(C.size_t(100)))
	if t.s == nil {
		//if t.s == (*C.char)(unsafe.Pointer(uintptr(0))) {
		panic("ban")
	}

	// 释放内存
	defer C.free(unsafe.Pointer(t.s))
	// 将go的字符串转为c的字符串，并自动释放
	ch = C.CString(s)
	defer C.free(unsafe.Pointer(ch))

	// 调用C的strncpy函数复制
	C.strncpy(t.s, ch, C.size_t(len(s)))
	// C的指针操作
	for i := C.size_t(0); i < C.strlen(t.s); i++ {
		tmp = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(t.s)) + uintptr(i)))
		*tmp = C.char(C.int(*tmp))
	}

	fmt.Println(C.GoString(t.s))
	fmt.Println(unsafe.Sizeof(t))
}
