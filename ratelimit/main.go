package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/kavu/go_reuseport"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

func handler(w http.ResponseWriter, req *http.Request, rateLimiter *redis_rate.Limiter) {
	userID := "user-12345"
	limit := int64(5)

	rate, reset, allowed := rateLimiter.AllowMinute(userID, limit)
	if !allowed {
		w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset, 10))
		http.Error(w, "API rate limit exceeded.", 429)
		return
	}

	fmt.Fprintf(w, "Hello world!\n")
	fmt.Fprint(w, "Rate limit remaining: ", strconv.FormatInt(limit-rate, 10))
}

func fasthandler(reqCtx *fasthttp.RequestCtx) {
	userID := "user-12345"
	limit := int64(5)

	rate, reset, allowed := mylimiter.AllowMinute(userID, limit)
	if !allowed {
		reqCtx.Response.Header.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		reqCtx.Response.Header.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
		reqCtx.Response.Header.Set("X-RateLimit-Reset", strconv.FormatInt(reset, 10))
		reqCtx.Error("API rate limit exceeded.", 429)
		return
	}

	reqCtx.Response.SetBodyString("Hello world!\n")
	reqCtx.Response.SetBodyString("Rate limit remaining: " + strconv.FormatInt(limit-rate, 10))
	return
}

func statusHandler(w http.ResponseWriter, req *http.Request, rateLimiter *redis_rate.Limiter) {
	userID := "user-12345"
	limit := int64(5)

	// With n=0 we just retrieve the current limit.
	rate, reset, allowed := rateLimiter.AllowN(userID, limit, time.Minute, 0)
	fmt.Fprintf(w, "Current rate: %d", rate)
	fmt.Fprintf(w, "Reset: %d", reset)
	fmt.Fprintf(w, "Allowed: %v", allowed)
}

var (
	ring      *redis.Ring
	mylimiter *redis_rate.Limiter
)

func init() {
	ring = redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379",
		},
	})
	mylimiter = redis_rate.NewLimiter(ring)
	// Optional.
	mylimiter.Fallback = rate.NewLimiter(rate.Every(time.Second), 100)
}

/**
func http() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		handler(w, req, mylimiter)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		statusHandler(w, req, mylimiter)
	})

	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("listening on localhost:8888...")
	log.Println(http.ListenAndServe("localhost:8888", nil))
	return
}
*/
func main() {
	router := fasthttprouter.New()
	router.GET("/", fasthandler)

	// start server

	listener, err := reuseport.Listen("tcp", ":8888")
	if err != nil {
		panic("error ")
		//logger.Fatal("server do not support reuse-port", zap.Error(err))
	}
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler: router.Handler,
		}
		if err := s.Serve(listener); err != nil {
			//	logger.Fatal("error in fast http server start", zap.Error(err))
			panic("error ")
		}

	}()

	//	logger.Info("aaa server start success  ", zap.String("server Address", config.Server.Addr))
	select {}

}
