package main

import (
	"github.com/tsingson/fastweb/fasthttputils"
	"github.com/valyala/fasthttp"
	"v/github.com/sanity-io/litter@v1.1.0"

	"io"
	"os"
)

func uploadHandler(ctx *fasthttp.RequestCtx) {
	path, err := fasthttputils.GetCurrentExecDir()
	if err != nil {
		return
	}
	litter.Dump(path)
	f, _ := os.OpenFile("./uploaded/"+string(ctx.FormValue("name")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	fheader, _ := ctx.FormFile("file")
	file, _ := fheader.Open()
	io.Copy(f, file)
	defer file.Close()
	defer f.Close()
}
