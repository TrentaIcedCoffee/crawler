package crawler

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type logger struct {
	error_stream  io.Writer
	output_stream io.Writer
}

var (
	DebugLog *log.Logger
)

func init() {
	DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
}

func (logger *logger) output(format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	if logger.output_stream == nil {
		fmt.Fprintf(os.Stdout, format, a...)
	} else {
		fmt.Fprintf(logger.output_stream, format, a...)
	}
}

func (logger *logger) error(format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	if logger.error_stream == nil {
		fmt.Fprintf(os.Stderr, format, a...)
	} else {
		fmt.Fprintf(logger.error_stream, format, a...)
	}
}

func (logger *logger) debug(format string, a ...any) {
	DebugLog.Printf(format, a...)
}
