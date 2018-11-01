package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	proxyAddr   string
	proxyClient = &fasthttp.HostClient{
		IsTLS: false,
		Addr:  "",
		// set other options here if required - most notably timeouts.
		// ReadTimeout: 60, // 如果在生产环境启用会出现多次请求现象
	}
)

func ReverseProxyHandler(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response

	prepareRequest(req)

	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
	}

	postprocessResponse(resp)
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(proxyAddr)
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func httpHandle(ctx *fasthttp.RequestCtx) {
	// 这里直接获取到 multipart.FileHeader, 需要手动打开文件句柄
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("get upload file error:", err)
		return
	}
	fh, err := f.Open()
	if err != nil {
		fmt.Println("open upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fh.Close() // 记得要关

	// 打开保存文件句柄
	fp, err := os.OpenFile("saveto.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open saving file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fp.Close() // 记得要关

	if _, err = io.Copy(fp, fh); err != nil {
		fmt.Println("save upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	ctx.Write([]byte("save file successfully!"))
}

func main() {
	port := flag.String("port", "8082", "listen port")
	targetAddr := flag.String("target", "api.yourdomain.com", "your server domain")
	flag.Parse()

	proxyClient.Addr = *targetAddr

	log.Println("port:", *port)
	log.Println("target:", *targetAddr)

	if err := fasthttp.ListenAndServe("localhost:"+*port, ReverseProxyHandler); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
}
