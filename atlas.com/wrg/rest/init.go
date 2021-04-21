package rest

import (
	"github.com/sirupsen/logrus"
)

func CreateRestService(l *logrus.Logger) {
	rs := NewServer(l)
	go rs.Run()
}
