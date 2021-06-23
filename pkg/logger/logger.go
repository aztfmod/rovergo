package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/console"
	"github.com/aztfmod/rover/pkg/utils"
	"github.com/sirupsen/logrus"
)

var stdOutLog = logrus.New()
var stdOutEntry *logrus.Entry
var fileLog = logrus.New()
var logLevel = logrus.InfoLevel
var currentFile string

type Fields = map[string]interface{}

// Only supporting fields on file logging
func WithFields(fields Fields) *logrus.Entry {
	stdOutEntry = stdOutLog.WithFields(fields)
	console.Infof("Level: %d", stdOutEntry.Level)
	return fileLog.WithFields(fields)
}

func Trace(args ...interface{}) {
	stdOutLog.Trace(args...)
	fileLog.Trace(args...)
}

func Debug(args ...interface{}) {
	stdOutLog.Debug(args...)
	fileLog.Debug(args...)
}

func Print(args ...interface{}) {
	stdOutLog.Print(args...)
	fileLog.Print(args...)
}

func Info(args ...interface{}) {
	if stdOutEntry == nil {
		console.Info("stdOutEntry is nil")
		stdOutLog.Info(args...)
	} else {
		console.Info("stdOutEntry is not nil")
		stdOutEntry.Info(args...)
	}
	fileLog.Info(args...)
}

func Warn(args ...interface{}) {
	stdOutLog.Warn(args...)
	fileLog.Warn(args...)
}

func Warning(args ...interface{}) {
	stdOutLog.Warning(args...)
	fileLog.Warning(args...)
}

func Error(args ...interface{}) {
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Error(args...)
	fileLog.Error(args...)
}

func Panic(args ...interface{}) {
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Panic(args...)
	fileLog.Panic(args...)
}

func Fatal(args ...interface{}) {
	//stdOutLog.Fatal(args...) - fatal calls Exit(1)
	stdOutLog.WithFields(logrus.Fields{"Fatal error encounterewd. For details see": currentFile}).Error(args...)
	fileLog.Fatal(args...) // So logging to file only
}

func init() {
	stdOutLog.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableColors:    false,
		FullTimestamp:    false,
		DisableTimestamp: true,
	})
	stdOutLog.SetLevel(logrus.InfoLevel)
	stdOutLog.SetOutput(os.Stdout)

	fileLog.SetLevel(logLevel)
	fileLog.SetFormatter(&logrus.JSONFormatter{})
	SetCommand("rover")
}

func SetLogLevel(level string, err error) {
	lowerLevel := strings.ToLower(level)
	switch lowerLevel {
	case "trace":
		logLevel = logrus.TraceLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "panic":
		logLevel = logrus.PanicLevel
	default:
		logLevel = logrus.InfoLevel
	}
	fileLog.SetLevel(logLevel)
}

func SetCommand(rovercommand string) {

	rovercommand = strings.ReplaceAll(rovercommand, " ", "")
	if rovercommand == "" {
		rovercommand = "rover"
	}

	roverHomeDir, _ := utils.GetRoverDirectory()
	roverHomeLogsDir := filepath.Join(roverHomeDir, "logs")
	command.EnsureDirectory(roverHomeLogsDir)

	// Add datetime stamp to filename
	currentTime := time.Now()
	currentFile = rovercommand + "--" + currentTime.Format("2006-01-02--15-04") + ".log"

	roverLogFile := filepath.Join(roverHomeLogsDir, currentFile)
	file, err := os.OpenFile(roverLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		fileLog.SetOutput(file)
	} else {
		stdOutLog.Info("Failed to log to file, using default stderr")
	}
}
