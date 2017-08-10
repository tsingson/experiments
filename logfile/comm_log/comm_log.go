package comm_log

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "time"
)

var default_log *Zlog

type Zlog struct {
    zap_log *zap.SugaredLogger
    curr_log_level zap.AtomicLevel
}

func Init(program_name, log_level_str, log_path string) *Zlog {
    encoder_cfg := zapcore.EncoderConfig{
        TimeKey:        "T",
        LevelKey:       "L",
        NameKey:        "N",
        CallerKey:      "C",
        MessageKey:     "M",
        StacktraceKey:  "S",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     time_encoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    log_level := zap.NewAtomicLevel()
    log_level.UnmarshalText([]byte(log_level_str))

    custom_cfg := zap.Config{
        Level:            log_level,
        Development:      true,
        Encoding:         "console",
        EncoderConfig:    encoder_cfg,
        OutputPaths:      []string{log_path},
        ErrorOutputPaths: []string{"stderr"},
    }

    logger, _ := custom_cfg.Build()
    new_logger := logger.Named(program_name)
    sugar := new_logger.Sugar()

    default_log = &Zlog{
        zap_log: sugar,
        curr_log_level : log_level,
    }

    return default_log
}

func Sync()  {
    default_log.zap_log.Sync()
}

func Debug(msg string, keysAndValues ...interface{}) {
    default_log.zap_log.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{})  {
    default_log.zap_log.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
    default_log.zap_log.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
    default_log.zap_log.Errorw(msg, keysAndValues...)
}

func Set_log_level(log_level string) {
    default_log.curr_log_level.UnmarshalText([]byte(log_level))
}

func Get_log_level() string {
    return default_log.curr_log_level.String()
}

func time_encoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
