package console

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// DebugEnabled controls output of debug messages and the spinner
var DebugEnabled = false
var consoleSpinner *spinner.Spinner

func init() {
	// See https://github.com/briandowns/spinner#available-character-sets
	consoleSpinner = spinner.New(spinner.CharSets[37], 100*time.Millisecond)
}

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

// Error outputs strings plus newline in red
func Error(s string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", s)
}

// Errorf outputs formatted strings in red
func Errorf(f string, a ...interface{}) {
	fmt.Printf("\033[1;31m"+f+"\033[0m", a...)
}

// Warning outputs strings plus newline in yellow
func Warning(s string) {
	fmt.Printf("\033[1;33m%s\033[0m\n", s)
}

// Warningf outputs formatted strings in yellow
func Warningf(f string, a ...interface{}) {
	fmt.Printf("\033[1;33m"+f+"\033[0m", a...)
}

// Success outputs strings plus newline in green
func Success(s string) {
	fmt.Printf("\033[1;32m%s\033[0m\n", s)
}

// Successf outputs formatted strings in green
func Successf(f string, a ...interface{}) {
	fmt.Printf("\033[1;32m"+f+"\033[0m", a...)
}

// Printfer implements the tfexec.printfer interface to be used with tfexec SetLogger
type Printfer struct{}

func (p Printfer) Printf(f string, a ...interface{}) {
	if !DebugEnabled {
		return
	}
	fmt.Printf("\033[1;35m"+f+"\033[0m", a...)
}

// StartSpinner starts the spinner, which is disabled when debug is set
func StartSpinner() {
	if DebugEnabled {
		return
	}
	consoleSpinner.Start()
}

// StopSpinner stops the spinner
func StopSpinner() {
	consoleSpinner.Stop()
}
