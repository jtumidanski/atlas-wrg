package handler

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type EmptyEventCreator func() interface{}

type EventHandler func(logrus.FieldLogger, opentracing.Span, interface{})