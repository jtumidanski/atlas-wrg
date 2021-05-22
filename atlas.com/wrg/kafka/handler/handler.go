package handler

import "github.com/sirupsen/logrus"

type EmptyEventCreator func() interface{}

type EventHandler func(logrus.FieldLogger, interface{})
