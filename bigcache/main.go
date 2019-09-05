package main

import (
	"log"
	"strings"
	"time"

	"github.com/allegro/bigcache"
	"github.com/valyala/fasthttp"
)

const pathPrefix = "/cache/"

var cache, cacheInitError = bigcache.NewBigCache(bigcache.Config{
	Shards:             256,              // number of shards (must be a power of 2)
	LifeWindow:         10 * time.Minute, // time after which entry can be evicted
	MaxEntriesInWindow: 10000 * 10 * 60,  // rps * lifeWindow
	MaxEntrySize:       1024,             // max entry size in bytes, used only in initial memory allocation
	Verbose:            true,             // prints information about additional memory allocation
})

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())
	key := path[len(pathPrefix):]

	if method == "GET" && strings.HasPrefix(path, pathPrefix) {
		if value, err := cache.Get(key); err == nil {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBody(value)
		} else {
			ctx.Error(err.Error(), fasthttp.StatusNotFound)
		}
	} else if method == "PUT" && strings.HasPrefix(path, pathPrefix) {
		cache.Set(key, ctx.Request.Body())
		ctx.SetStatusCode(fasthttp.StatusOK)
	} else {
		ctx.Error("Invalid path or method. Use GET or PUT with path: "+pathPrefix, fasthttp.StatusNotFound)
	}
}

func main() {
	if cacheInitError != nil {
		log.Fatal(cacheInitError.Error())
	}
	log.Fatal(fasthttp.ListenAndServe(":8080", fastHTTPHandler))
}
