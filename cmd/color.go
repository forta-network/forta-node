package cmd

import "github.com/fatih/color"

func yellowBold(str string, args ...interface{}) {
	color.New(color.Bold, color.FgYellow).Printf(str, args...)
}

func greenBold(str string, args ...interface{}) {
	color.New(color.Bold, color.FgGreen).Printf(str, args...)
}
