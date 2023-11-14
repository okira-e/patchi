package utils

import (
	"os"

	"github.com/Okira-E/patchi/pkg/vars/colors"
)

func Abort(message string) {
	PrintInColor(colors.Red, message)
	os.Exit(1)
}
