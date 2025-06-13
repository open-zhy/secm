package screen

import (
	"github.com/fatih/color"
)

var (
	printf  = color.New().Printf
	red     = color.New(color.FgRed).Printf
	redbold = color.New(color.FgRed, color.Bold).Printf
	blue    = color.New(color.FgBlue).Printf
	green   = color.New(color.FgGreen).Printf
	cyan    = color.New(color.FgCyan).Printf
)

func Printf(format string, args ...interface{}) {
	// This function is used to print messages without any color formatting.
	// It can be useful for logging or debugging purposes.
	// The format and args are passed directly to fmt.Printf.
	// Example usage: Printf("Hello, %s!", "world")
	// This will print: Hello, world!
	printf(format, args...)
}

func Println(format string) {
	// This function is used to print messages without any color formatting.
	// It can be useful for logging or debugging purposes.
	// The format and args are passed directly to fmt.Printf.
	// Example usage: Printf("Hello, %s!", "world")
	// This will print: Hello, world!
	printf(format + "\n")
}

func Infof(format string, args ...interface{}) {
	cyan(format, args...)
}

func Successf(format string, args ...interface{}) {
	green(format, args...)
}

func Errorf(format string, args ...interface{}) {
	red(format, args...)
}

func Redf(format string, args ...interface{}) {
	red(format, args...)
}

func RedBoldf(format string, args ...interface{}) {
	redbold(format, args...)
}

func Bluef(format string, args ...interface{}) {
	blue(format, args...)
}
