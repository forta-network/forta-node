package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func toStderr(str string) {
	fmt.Fprintf(os.Stderr, str)
}

func yellowBold(str string, args ...interface{}) {
	toStderr(color.New(color.Bold, color.FgYellow).Sprintf(str, args...))
}

func greenBold(str string, args ...interface{}) {
	color.New(color.Bold, color.FgGreen).Printf(str, args...)
}

func redBold(str string, args ...interface{}) {
	toStderr(color.New(color.Bold, color.FgRed).Sprintf(str, args...))
}

func whiteBold(str string, args ...interface{}) {
	color.New(color.Bold, color.FgWhite).Printf(str, args...)
}
