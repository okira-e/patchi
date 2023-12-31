package utils

import (
	"fmt"
	"os"

	"github.com/Okira-E/patchi/pkg/vars/colors"
)

// PrintInColor prints a string in the specified color.
// It takes a color code in ANSI format & a string to print.
// It also prints to stderr if the isError flag is provided.
// Example:
// PrintInColor("\033[31m", "This is red text.")
func PrintInColor(color string, str string, isError bool) {
	if isError {
		fmt.Fprintln(os.Stderr, color + str + colors.Reset)
	} else {
		fmt.Println(color + str + colors.Reset)
	}
}
