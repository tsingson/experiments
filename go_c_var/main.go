package main

/*
#include <stdio.h>
#include <stdlib.h>
char ch = 'M';
unsigned char uch = 253;
short st = 233;
int i = 257;
long lt = 11112222;
float f = 3.14;
double db = 3.15;
void * p;
char *str = "const string";
char str1[64] = "char array";
void printI(void *i) {
    printf("print i = %d\n", (*(int *)i));
}
static uint8_t des_key[] = { 0xa3, 0xbe, 0x93, 0xff, 0x10, 0x34, 0x5f, 0xde,
    0xc6, 0x2e, 0x57, 0x83, 0x29, 0x7c, 0x8e, 0xf6,
    0xa3, 0x58, 0x34, 0x27, 0x13, 0x2c, 0x4e, 0xd2 };
struct ImgInfo {
    char *imgPath;
    int format;
    unsigned int width;
    unsigned int height;
};
void printStruct(struct ImgInfo *imgInfo) {
    if(!imgInfo) {
        fprintf(stderr, "imgInfo is null\n");
        return ;
    }
    fprintf(stdout, "imgPath = %s\n", imgInfo->imgPath);
    fprintf(stdout, "format = %d\n", imgInfo->format);
    fprintf(stdout, "width = %d\n", imgInfo->width);
}

*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/sanity-io/litter"
)

func main() {
	fmt.Println("----------------Go to C---------------")
	fmt.Println(C.char('Y'))
	fmt.Printf("%c\n", C.char('Y'))
	fmt.Println(C.uchar('C'))
	fmt.Println(C.short(254))
	fmt.Println(C.long(11112222))
	var goi int = 2
	// unsafe.Pointer --> void *
	cpi := unsafe.Pointer(&goi)
	C.printI(cpi)
	fmt.Println("----------------C to Go---------------")
	fmt.Println(C.ch)
	fmt.Println(C.uch)
	fmt.Println(C.st)
	fmt.Println(C.i)
	fmt.Println(C.lt)
	f := float32(C.f)
	fmt.Println(reflect.TypeOf(f))
	fmt.Println(C.f)
	db := float64(C.db)
	fmt.Println(reflect.TypeOf(db))
	fmt.Println(C.db)
	// 区别常量字符串和char数组，转换成Go类型不一样
	str := C.GoString(C.str)
	fmt.Println(str)

	fmt.Println(reflect.TypeOf(C.str1))
	var charray []byte
	for i := range C.str1 {
		if C.str1[i] != 0 {
			charray = append(charray, byte(C.str1[i]))
		}
	}

	fmt.Println(charray)
	fmt.Println(string(charray))

	for i := 0; i < 10; i++ {
		imgInfo := C.struct_ImgInfo{imgPath: C.CString("../images/xx.jpg"), format: 0, width: 500, height: 400}
		defer C.free(unsafe.Pointer(imgInfo.imgPath))
		C.printStruct(&imgInfo)
	}

	fmt.Println("----------------C Print----------------")
	/**
	data := make([]byte, 1024)
	recved, _ := C.recv(C.int(5), unsafe.Pointer(&data), C.size_t(cap(data)), 0)
	litter.Dump(recved)
	*/
	var input_data = [1024]uint8{0xa3, 0xbe, 0x93, 0xff, 0x10, 0x34, 0x5f, 0xde,
		0xc6, 0x2e, 0x57, 0x83, 0x29, 0x7c, 0x8e, 0xf6,
		0xa3, 0x58, 0x34, 0x27, 0x13, 0x2c, 0x4e, 0xd2}
	litter.Dump(input_data)
	var input [1024]*C.uint8_t
	litter.Dump(unsafe.Pointer(&input))
}
