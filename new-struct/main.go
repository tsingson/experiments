package main

import (
	"fmt"
)

type options struct {
	a int64
	b string
	c map[int]string
}

func NewOption(opt ...ServerOption) *options {
	r := new(options)
	for _, o := range opt {
		o(r)
	}
	return r
}

type ServerOption func(*options)

func WriteA(s int64) ServerOption {
	return func(o *options) {
		o.a = s
	}
}

func WriteB(s string) ServerOption {
	return func(o *options) {
		o.b = s
	}
}

func WriteC(s map[int]string) ServerOption {
	return func(o *options) {
		o.c = s
	}
}

func main() {
	opt1 := WriteA(int64(1))
	opt2 := WriteB("test")
	opt3 := WriteC(make(map[int]string, 0))

	op := NewOption(opt1, opt2, opt3)

	fmt.Println(op.a, op.b, op.c)
}
