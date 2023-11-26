package utils

import (
	"github.com/gizak/termui/v3"
	"os"

	"github.com/Okira-E/patchi/pkg/vars/colors"
)

func Abort(message string) {
	PrintInColor(colors.Red, message)
	os.Exit(1)
}

func AbortTui(message string) {
	termui.Close()

	PrintInColor(colors.Red, message)
	os.Exit(1)
}
