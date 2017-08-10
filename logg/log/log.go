package log

import (
	"context"
	"os"
	"syscall"

	"github.com/pborman/uuid"
	"go.uber.org/zap"
)

type loggerKeyType struct{}

var (

	// G is a shorcut for GetLogger
	G = GetLogger

	// L is a global default logger
	L = zap.New(zap.NewTextEncoder(), zap.DebugLevel, zap.Output(os.Stderr))

	loggerKey = loggerKeyType{}
)

// WithLogger attached logger to a given context. Later the logger can be
// obtained by GetLogger
func WithLogger(ctx context.Context, logger zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// GetLogger either returns an attached logger from the context
// or global logger if nothing is attached
func GetLogger(ctx context.Context) zap.Logger {
	l := ctx.Value(loggerKey)
	if l == nil {
		return L
	}

	// NOTE: loggerKey is not accessable out of this package
	// so there only value assigned to that key is zap.Logger
	return l.(zap.Logger)
}

func genInstanceID() string {
	// NOTE: also pid info can be included
	hostname, err := os.Hostname()
	if err != nil {
		L.Warn("unable to get Hostname", zap.Error(err))
		hostname = uuid.NewRandom().String()
	}

	return hostname
}

func dup(f *os.File, name string) *os.File {
	fd, err := syscall.Dup(int(f.Fd()))
	if err != nil {
		L.Fatal("unable to dup fd", zap.String("name", name), zap.Error(err))
	}
	return os.NewFile(uintptr(fd), name)
}

func BuildOutput(filename string) *os.File {
	var output *os.File
	switch filename {
	case "", "/dev/stderr":
		output = dup(os.Stderr, "/dev/stderr")
	case "/dev/stdout":
		output = dup(os.Stdout, "/dev/stdout")
	default:
		var err error
		output, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			L.Fatal("unable to open logging output file", zap.Error(err))
		}
	}
	return output
}

func BuildContext(w zap.WriteSyncer, level zap.Level) context.Context {
	lvl := zap.DynamicLevel()
	lvl.SetLevel(level)
	outputOption := zap.Output(w)

	logger := zap.New(zap.NewTextEncoder(), outputOption, lvl)
	logger = logger.With(zap.String("instance", genInstanceID()))
	return WithLogger(context.Background(), logger)
}
