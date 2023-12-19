package utils

import (
	"fmt"
	"os"

	"github.com/gizak/termui/v3"

	"github.com/Okira-E/patchi/pkg/vars/colors"
)

// Abort stops execution and logs an error message to stderr in red.
func Abort(message string) {
	PrintInColor(colors.Red, message, true)
	os.Exit(1)
}

// AbortTui stops termui from rendering with a message to stderr.
func AbortTui(message string) {
	termui.Close()

	PrintInColor(colors.Red, message, true)
	os.Exit(1)
}

// LogTui stops termui from rendering with a message to be logged.
func LogTui(msg ...any) {
	termui.Close()

	fmt.Println(msg...)
	os.Exit(0)
}
