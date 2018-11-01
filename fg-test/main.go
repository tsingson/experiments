package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/koding/multiconfig"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsingson/fastweb/fasthttputils"
	"github.com/tsingson/fg"
	"github.com/tsingson/phi"
	"github.com/tsingson/vk/cmd/vklib/configuration"
	"github.com/valyala/fasthttp"
)

const (
	ServerName            = "Img service"
	Version               = "0.1.1-20180418"
	Concurrency           = 256 * 1024
	ReadBufferSize        = 1024 * 2
	MaxConnsPerIP         = 10
	MaxRequestsPerConn    = 10
	MaxRequestBodySize    = 100 << 20 // 100MB // 1024 * 4, // MaxRequestBodySize:
	VERSION               = "v1.0.2"
	LogFileNameTimeFormat = "2006-01-02-15"
	ConfigFileName        = "vkmsa-config.json"
	//	LogFileNameTimeFormat = "2006-01-02-15"
)

var (
	addr                                  = flag.String("addr", ":4243", "TCP address to listen to")
	compress                              = flag.Bool("compress", false, "Whether to enable transparent response compression")
	config                                *configuration.Config
	log                                   zerolog.Logger
	path, currentPath, scanPath, basePath string
)

func init() {
	afs := afero.NewOsFs()
	var err error
	{ // 获取当前目录
		path, err = fasthttputils.GetCurrentExecDir()
		if err != nil {
			panic("无法读取可执行程序的存储路径")
		}
		currentPath, err = fasthttputils.GetCurrentPath()
		if err != nil {
			panic("无法读取当前执行路径")
		}
	}
	{ // 读取配置文件
		configFile := path + "/" + ConfigFileName
		//	m := multiconfig.NewWithPath(configFile) // supports TOML and JSON
		m := &multiconfig.JSONLoader{Path: configFile}
		config = new(configuration.Config)
		//m.MustLoad(config) // Check for error
		err = m.Load(config)
		if err != nil {
			fmt.Println("config load fail")
			os.Exit(1)
		}
	}
	{
		// log setup
		logPath := path + "/log"
		logfileprefix := "vkmsa"

		check, _ := afero.DirExists(afs, logPath)
		if !check {
			afs.MkdirAll(logPath, 0755)
		}

		log = fasthttputils.NewZeroLog(logPath, logfileprefix)
	}

}

func requestHandler(ctx *fasthttp.RequestCtx) {
	begin := time.Now()
	time.Sleep(500 * time.Microsecond)
	fmt.Fprintf(ctx, "Hellx------------o~")
	end := time.Now()
	log.Info().Str("remoteIp", ctx.RemoteIP().String()).
		Str("agent", string(ctx.UserAgent())).
		Str("method", string(ctx.Method())).
		Str("url", string(ctx.RequestURI())).
		Int("status", ctx.Response.Header.StatusCode()).
		Time("start", begin).
		Dur("duration", end.Sub(begin)).Msg("fasthttp")
}

type FgLog struct {
	zerolog.Logger
	Log fg.Log
}

func main() {

	runtime.GOMAXPROCS(128)
	//	gspt.SetProcTitle("vkmsa")

	flag.Parse()

	var router *phi.Mux
	router = phi.NewRouter()

	httpLogger := fasthttputils.FastHttpZeroLogHandler
	router.Use(httpLogger)
	//
	router.Use(fasthttputils.Recoverer)
	router.Get("/", requestHandler)
	// grace fasthttp
	fglog := new(FgLog)
	fglog.Logger = log

	fastConfig := fg.FastConfig{
		Name:               ServerName,
		ReadBufferSize:     ReadBufferSize,
		MaxConnsPerIP:      MaxConnsPerIP,
		MaxRequestsPerConn: MaxRequestsPerConn,
		MaxRequestBodySize: MaxRequestBodySize, //  100 << 20, // 100MB // 1024 * 4, // MaxRequestBodySize:
		Concurrency:        Concurrency,
		FastLogger:         fasthttputils.NewFastHTTPLoggerAdapter(log),
		Logger:             fglog,
	}

	if err := fg.ListenAndServeConfig(*addr, router.ServeFastHTTP, fastConfig); err != nil {
		// fast http
		//if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatal().Err(err).Msg("Error in ListenAndServe")
	}
}
