package crawler

import (
	"fmt"
	"io"
	"os"
)

type logger struct {
	error_stream  io.Writer
	output_stream io.Writer
}

func (logger *logger) output(format string, a ...any) {
	format = format + "\n"
	if logger.output_stream == nil {
		fmt.Fprintf(os.Stdout, format, a...)
	} else {
		fmt.Fprintf(logger.output_stream, format, a...)
	}
}

func (logger *logger) error(format string, a ...any) {
	format = format + "\n"
	if logger.error_stream == nil {
		fmt.Fprintf(os.Stderr, format, a...)
	} else {
		fmt.Fprintf(logger.error_stream, format, a...)
	}
}

func (logger *logger) debug(format string, a ...any) {
	format = "[DEBUG] " + format + "\n"
	fmt.Fprintf(os.Stdout, format, a...)
}
