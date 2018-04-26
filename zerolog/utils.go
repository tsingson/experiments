package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	LogTimeFormat = "2006-01-02 15:04:05"
)

func time_encoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	enc.AppendString(t.Format(LogTimeFormat))
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

func InitLogFile(logFilename string) *lumberjack.Logger {
	lumberLogger := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // 开发时不压缩
	}
	//
	return lumberLogger
}
func InitLogger(encodeAsJSON, productMode bool, lumberLogger *lumberjack.Logger) *zap.Logger {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	// setup write or stdout

	logWriteSyncer := zapcore.Lock(zapcore.AddSync(lumberLogger))
	writersync := []zapcore.WriteSyncer{logWriteSyncer}
	if !productMode {
		writersync = append(writersync, os.Stdout)
	}
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     time_encoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encCfg)
	if encodeAsJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}
	atom := zap.NewAtomicLevel()

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writersync...), //logWriteSyncer,
		atom,
	)
	//zap.NewDevelopment()
	/**
	if productMode {
	 atom.SetLevel(zap.ErrorLevel)
	}
	*/
	logger := zap.New(core)
	//	defer logger.Sync()
	logger.Info("logger init success", zap.Time("InitTime", time.Now()))
	return logger
}

// design and code by tsingson
