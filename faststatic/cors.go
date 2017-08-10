package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func init() {

}

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func CORS(next fasthttprouter.Handler) fasthttprouter.Handle {
	return fasthttprouter.Handle(func(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx, p)
	})
}

/**

func Index(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
    fmt.Fprint(ctx, "some-api")
}

func main() {
    var handler fasthttprouter.Handle
    handler = Index
    handler = CORS(handler)

    router.Handle("GET", "/", handler)
}


*/

// design and code by tsingson
