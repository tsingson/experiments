package main

import (
	"log"
	"os"

	"time"

	"github.com/valyala/fasthttp"
)

var (
	output = log.New(os.Stdout, "", 0)
)

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// Tiny format:
// <method> <url> - <status> - <response-time us>
// GET / - 200 - 11.925 us
func Tiny(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%s %s - %v - %v",
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	})
}

// Short format:
// <remote-addr> | <HTTP/:http-version> | <method> <url> - <status> - <response-time us>
// 127.0.0.1:53324 | HTTP/1.1 | GET /hello - 200 - 44.8µs
func Short(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%v | %s | %s %s - %v - %v",
			ctx.RemoteAddr(),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	})
}

// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url> - <status> - <response-time us> | <user-agent>
// [2017/05/31 - 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello - 200 - 48.279µs | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func Combined(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}
