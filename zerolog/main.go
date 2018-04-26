// +build !binary_log

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logFilename := "./logxxxxxx.log"
	maxSize := 10
	lumberLogger := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    maxSize, // megabytes
		MaxBackups: 31,
		MaxAge:     31,    // days
		Compress:   false, // 开发时不压缩
	}

	d := diodes.NewManyToOne(1024*1024*2, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	w := diode.NewWriter(lumberLogger, d, 10*time.Millisecond)
	zerolog.TimeFieldFormat = ""
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	log := zerolog.New(w)
	for i := 0; i < 1024*1024*1024; i++ {
		log.Debug().
			Str("Scale", "833 cents").
			Float64("Interval", 833.09).
			Msg("Fibonacci is everywhere_test  " + strconv.Itoa(i))
	}
	defer w.Close()
}
func _main() {
	log := initLogger()

	log.Print("hello world")
	//logger := fastlog.New(os.Stderr).With().Timestamp().Logger()
	log.Output(os.Stderr)
	log.Info().Str("foo", "bar").Msg("hello world")

}

func initLogger() zerolog.Logger {
	logFilename := "./log.log"
	maxSize := 1
	lumberLogger := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    maxSize, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // 开发时不压缩
	}
	d := diodes.NewManyToOne(10000, diodes.AlertFunc(func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	}))
	w := diode.NewWriter(lumberLogger, d, 100*time.Millisecond)
	//	w := diode.NewWriter(os.Stdout, d, 10*time.Millisecond)
	zerolog.TimeFieldFormat = ""
	log := zerolog.New(w)
	/**
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			lumberLogger.Rotate()
		}
	}()
	*/
	return log
	// Output: {"level":"debug","message":"test"}
}
