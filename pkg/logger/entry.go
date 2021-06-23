package logger

import "github.com/sirupsen/logrus"

type Entry Fields

func (logFields *Entry) Error(args ...interface{}) {
	f := logrus.Fields(*logFields)
	fileEntry := fileLog.WithFields(f)
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Error(args...)
	fileEntry.Error(args...)
}

func (logFields *Entry) Info(message string) {
	f := logrus.Fields(*logFields)
	fileEntry := fileLog.WithFields(f)
	stdOutLog.Info(message)
	fileEntry.Info(message)
}
