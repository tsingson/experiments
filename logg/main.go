package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	logfile     string = "./output.log"
	errorfile   string = "./error.log"
	url         string = "test localhost"
	encoder_cfg        = zapcore.EncoderConfig{
		/**
		TimeKey:       "T",

		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		*/
		TimeKey:       "eventTime",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",

		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     time_encoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

func time_encoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func ZapProductionConfig(logfile, errorfile string) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		//Encoding: "json",
		Encoding: "console",
		//EncoderConfig:   zap.NewProductionEncoderConfig(),
		EncoderConfig:    encoder_cfg,
		OutputPaths:      []string{"stderr", logfile},
		ErrorOutputPaths: []string{"stderr", errorfile},
	}
}
func main() {

	logger, _ := ZapProductionConfig(logfile, errorfile).Build()

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Debugw("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infow("Failed   URL: %s", url, "failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second)
	logger.Sync()

}
