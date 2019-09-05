package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
}

func logInFile() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}

func sigHup(lumberLogger *lumberjack.Logger) {
	//	lubberLogger := &lumberjack.Logger{}
	//	log.SetOutput(lubberLogger)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			lumberLogger.Rotate()
		}
	}()
}

func InitLogger() *zap.Logger {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	lumberLogger := &lumberjack.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    // days
		Compress:   false, // 开发时不压缩
	}
	//
	sigHup(lumberLogger)

	logWriteSyncer := zapcore.AddSync(lumberLogger)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		logWriteSyncer,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	return logger
}

// design and code by tsingson
