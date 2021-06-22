package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aztfmod/rover/pkg/command"
	"github.com/aztfmod/rover/pkg/utils"
	log "github.com/sirupsen/logrus"
)

var LogEnvironment = "production"
var logLevel = log.InfoLevel

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableColors:    false,
		FullTimestamp:    false,
		DisableTimestamp: true,
	})

	log.SetLevel(logLevel)

	if LogEnvironment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		roverHomeDir, _ := utils.GetRoverDirectory()
		roverHomeLogsDir := filepath.Join(roverHomeDir, "logs")
		command.EnsureDirectory(roverHomeLogsDir)
		roverLogFile := filepath.Join(roverHomeLogsDir, "rover.log")
		file, err := os.OpenFile(roverLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	} else {
		log.SetReportCaller(true) // has performance impact but useful in debugging
	}

}

func SetLogLevel(level string, err error) {
	lowerLevel := strings.ToLower(level)
	switch lowerLevel {
	case "trace":
		logLevel = log.TraceLevel
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "fatal":
		logLevel = log.FatalLevel
	case "panic":
		logLevel = log.PanicLevel
	default:
		logLevel = log.InfoLevel
	}
}
