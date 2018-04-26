package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type (
	// FastHTTPLoggerAdapter  Adapter for passing apex/log used as gramework Logger into fasthttp
	FastHTTPLoggerAdapter struct {
		ZeroLogger zerolog.Logger
		fasthttp.Logger
	}
)

var (
	log zerolog.Logger
	err error
)

func init() {
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func main() {
	/**
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	req.Header.Add("Accept-Encoding", "gzip")
	req.SetRequestURI("http://localhost:8080")

	err := fasthttp.HostClient{}(req, res)
	if err != nil {
		panic(err)
	}
	*/
	server := fasthttp.Server{
		Name:               "Fasthttp server",
		Handler:            handler,
		ReduceMemoryUsage:  true,
		ReadBufferSize:     16000,
		WriteBufferSize:    16000,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 10,
		MaxRequestBodySize: 1024 * 4, // MaxRequestBodySize: 100<<20, // 100MB

		Logger: NewFastHTTPLoggerAdapter(log),
	}
	err = server.ListenAndServe(":8080")
	if err != nil {
		log.Panic().Err(err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")
	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
	ctx.SetContentType("text/plain; charset=utf8")
}

// NewFastHTTPLoggerAdapter create new *FastHTTPLoggerAdapter
func NewFastHTTPLoggerAdapter(logger zerolog.Logger) (fasthttplogger *FastHTTPLoggerAdapter) {
	fasthttplogger = &FastHTTPLoggerAdapter{
		ZeroLogger: logger,
	}
	return fasthttplogger
}
