package kafka

import (
	"atlas-wrg/retry"
	"atlas-wrg/topic"
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

func CreateConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, configs ...ConsumerConfig) {
	for _, c := range configs {
		go createConsumer(l, ctx, wg, c)
	}
}

func NewConsumerConfig[E any](name string, topicToken string, groupId string, handler HandlerFunc[E]) ConsumerConfig {
	return ConsumerConfig{
		name:       name,
		topicToken: topicToken,
		groupId:    groupId,
		maxWait:    500,
		handler:    adapt(handler),
	}
}

type ConsumerConfig struct {
	name       string
	topicToken string
	groupId    string
	maxWait    time.Duration
	handler    messageHandler
}

type messageHandler func(l logrus.FieldLogger, span opentracing.Span, msg kafka.Message)

type HandlerFunc[E any] func(logrus.FieldLogger, opentracing.Span, E)

func adapt[E any](eh HandlerFunc[E]) messageHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, msg kafka.Message) {
		var event E
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			l.WithError(err).Errorf("Could not unmarshal event into %s.", msg.Value)
		} else {
			eh(l, span, event)
		}
	}
}

func createConsumer(cl *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, c ConsumerConfig) {
	initSpan := opentracing.StartSpan("consumer_init")
	t := topic.GetRegistry().Get(cl, initSpan, c.topicToken)
	initSpan.Finish()

	l := cl.WithFields(logrus.Fields{"originator": t, "type": "kafka_consumer"})

	l.Infof("Creating topic consumer.")

	wg.Add(1)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   t,
		GroupID: c.groupId,
		MaxWait: c.maxWait,
	})

	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		for {
			var msg kafka.Message
			readerFunc := func(attempt int) (bool, error) {
				var err error
				msg, err = r.ReadMessage(ctx)
				if err == io.EOF || err == context.Canceled {
					return false, err
				} else if err != nil {
					l.WithError(err).Warnf("Could not read message on topic %s, will retry.", r.Config().Topic)
					return true, err
				}
				return false, err
			}

			err := retry.Try(readerFunc, 10)
			if err == io.EOF || err == context.Canceled || len(msg.Value) == 0 {
				l.Infof("Reader closed, shutdown.")
				return
			} else if err != nil {
				l.WithError(err).Errorf("Could not successfully read message.")
			} else {
				l.Infof("Message received %s.", string(msg.Value))
				go func() {
					headers := make(map[string]string)
					for _, header := range msg.Headers {
						headers[header.Key] = string(header.Value)
					}

					spanContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
					span := opentracing.StartSpan(c.name, opentracing.FollowsFrom(spanContext))
					defer span.Finish()

					c.handler(l, span, msg)
				}()
			}
		}
	}()

	l.Infof("Start consuming topic.")
	<-ctx.Done()
	l.Infof("Shutting down topic consumer.")
	if err := r.Close(); err != nil {
		l.WithError(err).Errorf("Error closing reader.")
	}
	wg.Done()
	l.Infof("Topic consumer stopped.")
}
