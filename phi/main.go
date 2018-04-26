package main

import (
	"fmt"
	"log"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	//	output = log.New(os.Stdout, "", 0)
	output zerolog.Logger
)

func init() {
	logFilename := "./logxxxxxx.log"
	maxSize := 1
	lumberLogger := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    maxSize, // megabytes
		MaxBackups: 31,
		MaxAge:     31,    // days
		Compress:   false, // 开发时不压缩
	}

	d := diodes.NewManyToOne(1024*1024*4, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	w := diode.NewWriter(lumberLogger, d, 100*time.Millisecond)
	zerolog.TimeFieldFormat = ""
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	output = zerolog.New(w)
	//	defer w.Close()
}

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

func main() {
	r := phi.NewRouter()

	httpLogger := func(next phi.HandlerFunc) phi.HandlerFunc {
		return func(ctx *fasthttp.RequestCtx) {
			begin := time.Now()
			next(ctx)
			end := time.Now()
			output.Info().
				Str("remoteAddr", ctx.RemoteAddr().String()).
				Str("remoteIp", ctx.RemoteIP().String()).
				Str("agent", string(ctx.UserAgent())).
				Str("method", string(ctx.Method())).
				Str("url", string(ctx.RequestURI())).
				Int("status", ctx.Response.Header.StatusCode()).
				Time("start", begin).
				Dur("duration", end.Sub(begin)).
				Time("end", end).
				Msg("")
			/**
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
			*/
		}
	}
	r.Use(httpLogger)
	/**
	reqIDMW := func(next phi.HandlerFunc) phi.HandlerFunc {
		return func(ctx *fasthttp.RequestCtx) {
			next(ctx)
			ctx.WriteString("+reqid=1")
		}
	}
	r.Use(reqIDMW)
	*/
	r.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("index")
	})

	r.NotFound(func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("whoops, not found")
		ctx.SetStatusCode(404)
	})
	r.MethodNotAllowed(func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("whoops, bad method")
		ctx.SetStatusCode(405)
	})

	// cat
	r.Route("/cat", func(r phi.Router) {
		r.NotFound(func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("no such cat")
			ctx.SetStatusCode(404)
		})
		r.Use(func(next phi.HandlerFunc) phi.HandlerFunc {
			return func(ctx *fasthttp.RequestCtx) {
				next(ctx)
				ctx.WriteString("+cat")
			}
		})
		r.Get("/cat", func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("cat")
		})
		r.Patch("/", func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("patch cat")
		})
	})

	// user
	userRouter := phi.NewRouter()
	userRouter.NotFound(func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("no such user")
		ctx.SetStatusCode(404)
	})
	userRouter.Use(func(next phi.HandlerFunc) phi.HandlerFunc {
		return func(ctx *fasthttp.RequestCtx) {
			next(ctx)
			ctx.WriteString("+user")
		}
	})
	userRouter.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("user")
	})
	userRouter.Post("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("new user")
	})
	r.Mount("/user", userRouter)

	server := &fasthttp.Server{
		Handler:       r.ServeFastHTTP,
		ReadTimeout:   10 * time.Second,
		MaxConnsPerIP: 1024 * 4,
	}

	log.Fatal(server.ListenAndServe(":7789"))
}
