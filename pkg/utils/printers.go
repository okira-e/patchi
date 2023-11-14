package utils

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/vars/colors"
)

// PrintInColor prints a string in the specified color.
// It takes a color code in ANSI format & a string to print.
// Example:
// PrintInColor("\033[31m", "This is red text.")
func PrintInColor(color string, str string) {
	fmt.Println(color + str + colors.Reset)
}
