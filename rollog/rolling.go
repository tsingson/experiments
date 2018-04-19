package main

import (
	"github.com/robfig/cron"
	"github.com/sanity-io/litter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

const (
	port = ":9999"
)

type (
	Config struct {
		// EncodeLogsAsJson makes the log framework log JSON
		EncodeLogsAsJson bool
		// FileLoggingEnabled makes the framework log to a file
		// the fields below can be skipped if this value is false!
		FileLoggingEnabled bool
		// Directory to log to to when filelogging is enabled
		Directory string
		// Filename is the name of the logfile which will be placed inside the directory
		Filename string
		// MaxSize the max size in MB of the logfile before it's rolled
		MaxSize int
		// MaxBackups the max number of rolled files to keep
		MaxBackups int
		// MaxAge the max age in days to keep a logfile
		MaxAge int
	}
)

func setupLogger() {
	config := Config{
		// EncodeLogsAsJson: false,
		EncodeLogsAsJson:   true,
		FileLoggingEnabled: true,
		Directory:          "./logs",
		Filename:           "servant.log",
		MaxAge:             30,
		MaxSize:            1000,
	}
	Configure(config)
}

// Configuration for logging

var DefaultZapLogger = newZapLogger(false, os.Stdout)

func Configure(config Config) {
	writers := []zapcore.WriteSyncer{os.Stdout}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	DefaultZapLogger = newZapLogger(config.EncodeLogsAsJson, zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(DefaultZapLogger)

}

func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, 0755); err != nil {
		//	Error("failed create log directory", zap.Error(err), zap.String("path", config.Directory))
		panic("logger init fail.")
		return nil
	}

	lj_log := lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    //megabytes
		MaxAge:     config.MaxAge,     //days
		MaxBackups: config.MaxBackups, //files
		LocalTime:  true,
	}

	c := cron.New()
	// c.AddFunc("* * * * * *", func() { lj_log.Rotate() })
	c.AddFunc("@daily", func() { lj_log.Rotate() })
	c.Start()

	return zapcore.AddSync(&lj_log)
}

func newZapLogger(encodeAsJSON bool, output zapcore.WriteSyncer) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "logtime",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	opts := []zap.Option{zap.AddCaller()}
	opts = append(opts, zap.AddStacktrace(zap.WarnLevel))
	encoder := zapcore.NewConsoleEncoder(encCfg)
	if encodeAsJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	return zap.New(zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(zap.DebugLevel)), opts...)
}

func main() {
	setupLogger()
	zaplogger := DefaultZapLogger
	// defer logger.Sync() // flushes buffer, if any
	litter.Dump(zaplogger)
	fmt.Println("running...")

}
