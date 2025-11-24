package logger

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// Info prints a message in Cyan
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[INFO] %s%s\n", Cyan, msg, Reset)
}

// Success prints a message in Green
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[SUCCESS] %s%s\n", Green, msg, Reset)
}

// Warning prints a message in Yellow
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[WARN] %s%s\n", Yellow, msg, Reset)
}

// Error prints a message in Red
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[ERROR] %s%s\n", Red, msg, Reset)
}

// Fatal prints a message in Red and exits
func Fatal(format string, args ...interface{}) {
	Error(format, args...)
	os.Exit(1)
}
