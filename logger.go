package crawler

import (
	"fmt"
	"os"
)

type logger struct {
	isDebug bool
}

func (logger *logger) Debug(format string, a ...any) {
	if logger.isDebug {
		fmt.Println("[DEBUG] " + fmt.Sprintf(format, a...))
	}
}

func (logger *logger) Output(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (logger *logger) Error(format string, a ...any) {
	fmt.Fprintln(os.Stderr, "[ERROR] "+fmt.Sprintf(format, a...))
}
