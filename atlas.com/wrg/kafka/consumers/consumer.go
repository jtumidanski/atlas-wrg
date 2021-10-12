package consumers

import (
	"atlas-wrg/kafka/handler"
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

type config struct {
	maxWait time.Duration
}

type ConfigOption func(c *config)

func NewConsumer(cl *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, groupId string, ec handler.EmptyEventCreator, h handler.EventHandler, modifications ...ConfigOption) {
	c := &config{maxWait: 500 * time.Millisecond}

	for _, modification := range modifications {
		modification(c)
	}

	initSpan := opentracing.StartSpan("consumer_init")
	t := topic.GetRegistry().Get(cl, initSpan, topicToken)
	initSpan.Finish()

	l := cl.WithFields(logrus.Fields{"originator": t, "type": "kafka_consumer"})

	l.Infof("Creating topic consumer.")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   t,
		GroupID: groupId,
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
				event := ec()
				err = json.Unmarshal(msg.Value, &event)
				if err != nil {
					l.WithError(err).Errorf("Could not unmarshal event into %s.", msg.Value)
				} else {
					go func() {
						headers := make(map[string]string)
						for _, header := range msg.Headers {
							headers[header.Key] = string(header.Value)
						}

						spanContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
						span := opentracing.StartSpan(name, opentracing.FollowsFrom(spanContext))
						defer span.Finish()

						h(l, span, event)
					}()
				}
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
