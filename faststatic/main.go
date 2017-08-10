// Example static file server.
//
// Serves static files from the given directory.
// Exports various stats at /stats .
package main

import (
	"expvar"
	"flag"
	"github.com/buaazp/fasthttprouter"
	"github.com/kavu/go_reuseport"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"

	"log"
	"net"
	"runtime"
)

var (
	addr               = flag.String("addr", ":8081", "TCP address to listen to")
	addrTLS            = flag.String("addrTLS", "", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS")
	byteRange          = flag.Bool("byteRange", false, "Enables byte range requests if set to true")
	certFile           = flag.String("certFile", "./ssl-cert-snakeoil.pem", "Path to TLS certificate file")
	compress           = flag.Bool("compress", true, "Enables transparent response compression if set to true")
	dir                = flag.String("dir", "/Users/qinshen/git/g2cn/public", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", false, "Whether to generate directory index pages")
	keyFile            = flag.String("keyFile", "./ssl-cert-snakeoil.key", "Path to TLS key file")
	vhost              = flag.Bool("vhost", false, "Enables virtual hosting by prepending the requested path with the requested hostname")
	listener           net.Listener
	slistener          net.Listener
	err                error
)

func main() {
	// Parse command-line flags.

	// Create RequestHandler serving server stats on /stats and files
	// on other requested paths.
	// /stats output may be filtered using regexps. For example:
	//
	//   * /stats?r=fs will show only stats (expvars) containing 'fs'
	//     in their names.
	/**
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stats":
			expvarhandler.ExpvarHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx)
		}
	}
	*/

	runtime.GOMAXPROCS(runtime.NumCPU())
	// 使用 reuseport 进行连接监听, 以便使用 taskset 来绑定 CPU

	//defer listener.Close()

	router := fasthttprouter.New()
	//	router.GET("/", fsHandler)
	//router.NotFound = fasthttp.FSHandler(*dir, 0)
	//
	router.NotFound = fsHandler(*dir, 0)

	router.GET("/stats/", expvarhandler.ExpvarHandler)
	//	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	// Start HTTP server.
	/*
		if len(*addr) > 0 {
			log.Printf("Starting HTTP server on %q", *addr)
			go func() {
				if err := fasthttp.Serve(listener, router.Handler); err != nil {
					log.Fatalf("error in ListenAndServe: %s", err)
				}
			}()
		}

		// Start HTTPS server.

			if len(*addrTLS) > 0 {
				log.Printf("Starting HTTPS server on %q", *addrTLS)
				go func() {
					if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, router.Handler); err != nil {
						log.Fatalf("error in ListenAndServeTLS: %s", err)
					}
				}()
			}
	*/

	if len(*addr) > 0 {
		listener, err = reuseport.Listen("tcp", *addr)
		if err != nil {
			panic(err)
		}
		go func() {
			if err := fasthttp.Serve(listener, router.Handler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}
	if len(*addrTLS) > 0 {
		log.Printf("Starting HTTPS server on %q", *addrTLS)
		slistener, err = reuseport.Listen("tcp", *addr)
		if err != nil {
			panic(err)
		}
		go func() {
			if err := fasthttp.ServeTLS(slistener, *certFile, *keyFile, router.Handler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
	}
	log.Printf("Serving files from directory %q", *dir)
	log.Printf("See stats at http://%s/stats", *addr)

	// Wait forever.
	select {}
}

func fsHandler(root string, stripSlashes int) fasthttp.RequestHandler {
	// Setup FS handler
	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		Compress:           *compress,
		AcceptByteRange:    *byteRange,
	}
	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewVHostPathRewriter(stripSlashes)
	}
	return fs.NewRequestHandler()
}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
