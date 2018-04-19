package main

import (
	"bufio"
	"io"
	"os"
)

func processBlock(line []byte) {
	os.Stdout.Write(line)
}

func ReadBlock(filePth string, bufSize int, hookfn func([]byte)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, bufSize) //一次读取多少个字节
	bfRd := bufio.NewReader(f)
	for {
		n, err := bfRd.Read(buf)
		hookfn(buf[:n]) // n 是成功读取字节数

		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

func main() {
	ReadBlock("test.txt", 10000, processBlock)
}
