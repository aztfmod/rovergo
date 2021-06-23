package logger

import "github.com/sirupsen/logrus"

type Entry Fields

func (logFields *Entry) Error(args ...interface{}) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Error(args...)
	fileEntry.Error(args...)
}

func (logFields *Entry) Fatal(args ...interface{}) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.WithFields(logrus.Fields{"Fatal error encounterewd. For details see": currentFile}).Error(args...)
	fileEntry.Fatal(args...)
}

func (logFields *Entry) Panic(args ...interface{}) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Panic(args...)
	fileEntry.Panic(args...)
}

func (logFields *Entry) Trace(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Trace(message)
	fileEntry.Trace(message)
}

func (logFields *Entry) Debug(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Debug(message)
	fileEntry.Debug(message)
}

func (logFields *Entry) Info(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Info(message)
	fileEntry.Info(message)
}

func (logFields *Entry) Warn(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Warn(message)
	fileEntry.Warn(message)
}

func (logFields *Entry) Warning(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Warning(message)
	fileEntry.Warning(message)
}
