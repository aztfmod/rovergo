package console

import (
	"fmt"
)

// DebugEnabled controls output of debug messages
var DebugEnabled = false

// Debug outputs strings if debug is enabled, in delightful shade of magenta
func Debug(s string) {
	if !DebugEnabled {
		return
	}
	fmt.Printf("\033[1;35m%s\033[0m\n", s)
}

// Debugf outputs formatted strings if debug is enabled, in delightful shade of magenta
func Debugf(f string, a ...interface{}) {
	if !DebugEnabled {
		return
	}
	fmt.Printf("\033[1;35m"+f+"\033[0m", a...)
}

// Info outputs strings plus newline in blue
func Info(s string) {
	fmt.Printf("\033[1;34m%s\033[0m\n", s)
}

// Infof outputs formatted strings in blue
func Infof(f string, a ...interface{}) {
	fmt.Printf("\033[1;34m"+f+"\033[0m", a...)
}

// Error outputs strings plus newline in blue
func Error(s string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", s)
}

// Infof outputs formatted strings in red
func Errorf(f string, a ...interface{}) {
	fmt.Printf("\033[1;31m"+f+"\033[0m", a...)
}

// Warning outputs strings plus newline in blue
func Warning(s string) {
	fmt.Printf("\033[1;33m%s\033[0m\n", s)
}

// Infof outputs formatted strings in yellow
func Warningf(f string, a ...interface{}) {
	fmt.Printf("\033[1;33m"+f+"\033[0m", a...)
}

// Success outputs strings plus newline in blue
func Success(s string) {
	fmt.Printf("\033[1;32m%s\033[0m\n", s)
}

// Infof outputs formatted strings in green
func Successf(f string, a ...interface{}) {
	fmt.Printf("\033[1;32m"+f+"\033[0m", a...)
}
