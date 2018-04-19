package main

/**
#include <stdio.h>
char* Calc();

int a=10000, b, c=2800, d, e, f[2801], g,i;
char r[1000];
char* pr = r;
char* Calc() {
	for(;b-c;)
	f[b++]=a/5;
	//for(;d=0,g=c*2;c-=14,printf("%.4d",e+d/a),e=d%a)
	for(;d=0,g=c*2;c-=14,sprintf(pr,"%.4d",e+d/a),pr +=4,e=d%a)
	for(b=c;d+=f[b]*a,f[b]=d%--g,d/=g--,--b;d*=b);
	return r;
}
*/
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println("计算PI值:")
	v := C.GoString(C.Calc())
	fmt.Println(v)
}
