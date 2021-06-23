package logger

import "github.com/sirupsen/logrus"

type Entry Fields

func (logFields *Entry) Error(args ...interface{}) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.WithFields(logrus.Fields{"For details see": currentFile}).Error(args...)
	fileEntry.Error(args...)
}

func (logFields *Entry) Info(message string) {
	fileEntry := fileLog.WithFields(logrus.Fields(*logFields))
	stdOutLog.Info(message)
	fileEntry.Info(message)
}
