package main

import (
	"crypto/md5"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"html/template"
	"io"
	"log"
	"strconv"
	"time"
)

// Hello Handler
func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

// Index Handler, 默认的ContenType是text/plain, 输出的内容在pre标签里面
// 因此这里必须手动设置ContentType为text/html
func Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("text/html")
	t := template.New("index.gtpl")
	t, _ = t.Parse(`<html>
    <head>
        <title>上传文件</title>
    </head>
    <body>
    <form enctype="multipart/form-data" action="/upload" method="post">
      <input type="file" name="uploadfile" />
      <input type="hidden" name="token" value="{{.}}"/>
      <input type="submit" value="upload" />
    </form>
    </body>
</html>`)
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	t.Execute(ctx, token)
}

// UploadHandler is here
func UploadHandler(ctx *fasthttp.RequestCtx) {
	data, err := ctx.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("get upload file error:", err)
		return
	}
	fileObj := data.File["uploadfile"][0]
	err = fasthttp.SaveMultipartFile(fileObj, fileObj.Filename)
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("save upload file error:", err)
		return
	}
	ctx.Write([]byte("save file successfully!"))
}

// DownloadHandler is here
func DownloadHandler(ctx *fasthttp.RequestCtx) {
	fileName := ctx.UserValue("filename")
	switch fileName := fileName.(type) {
	case string:
		ctx.SendFile(fileName)
	default:
		ctx.SetStatusCode(500)
		fmt.Println("the filename is not string.")
	}
}
func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.POST("/upload", UploadHandler)
	router.GET("/download/:filename", DownloadHandler)
	router.GET("/hello/:name", Hello)
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
