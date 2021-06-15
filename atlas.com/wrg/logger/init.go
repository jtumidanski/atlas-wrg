package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func CreateLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{})
	if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := logrus.ParseLevel(val); err == nil {
			l.SetLevel(level)
		}
	}
	return l
}
