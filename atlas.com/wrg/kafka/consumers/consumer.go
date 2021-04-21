package consumers

import (
	"atlas-wrg/kafka/topics"
	"atlas-wrg/retry"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Consumer struct {
	l                 log.FieldLogger
	ctx               context.Context
	name              string
	groupId           string
	topicToken        string
	emptyEventCreator EmptyEventCreator
	h                 EventProcessor
}

func NewConsumer(l *log.Logger, ctx context.Context, h EventProcessor, options ...ConsumerOption) Consumer {
	c := &Consumer{}
	c.ctx = ctx
	c.h = h
	for _, option := range options {
		option(c)
	}

	c.name = topics.GetRegistry().Get(l, c.topicToken)
	c.l = l.WithFields(log.Fields{"originator": c.name, "type": "kafka_consumer"})
	return *c
}

type EmptyEventCreator func() interface{}

type EventProcessor func(log.FieldLogger, interface{})

type ConsumerOption func(c *Consumer)

func SetGroupId(groupId string) func(c *Consumer) {
	return func(c *Consumer) {
		c.groupId = groupId
	}
}

func SetTopicToken(topicToken string) func(c *Consumer) {
	return func(c *Consumer) {
		c.topicToken = topicToken
	}
}

func SetEmptyEventCreator(f EmptyEventCreator) func(c *Consumer) {
	return func(c *Consumer) {
		c.emptyEventCreator = f
	}
}

func (c Consumer) Init() {
	c.l.Infof("Creating topic consumer.")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   c.name,
		GroupID: c.groupId,
		MaxWait: 500 * time.Millisecond,
	})

	for {
		msg, err := retry.RetryResponse(consumerReader(c.l, r, c.ctx), 10)
		if err != nil {
			c.l.WithError(err).Errorf("Could not successfully read message.")
		} else {
			if val, ok := msg.(*kafka.Message); ok {
				event := c.emptyEventCreator()
				err = json.Unmarshal(val.Value, &event)
				if err != nil {
					c.l.WithError(err).Errorf("Could not unmarshal event into %s.", val.Value)
				} else {
					c.h(c.l, event)
				}
			} else {
				c.l.Errorf("Message received not a valid kafka message.")
			}
		}
	}
}

func consumerReader(l log.FieldLogger, r *kafka.Reader, ctx context.Context) retry.RetryResponseFunc {
	return func(attempt int) (bool, interface{}, error) {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			l.WithError(err).Warnf("Could not read message on topic %s, will retry.", r.Config().Topic)
			return true, nil, err
		}
		return false, &msg, err
	}
}
