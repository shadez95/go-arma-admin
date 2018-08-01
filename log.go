package main

import "github.com/sirupsen/logrus"

func configureLogger(loglevel string) {
	switch loglevel {
	case "debug":
		Log.Level = logrus.DebugLevel
	case "info":
		Log.Level = logrus.InfoLevel
	case "warn":
		Log.Level = logrus.WarnLevel
	case "error":
		Log.Level = logrus.ErrorLevel
	case "fatal":
		Log.Level = logrus.FatalLevel
	case "panic":
		Log.Level = logrus.PanicLevel
	default:
		Log.Level = logrus.InfoLevel
	}
}
