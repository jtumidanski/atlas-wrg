package consumers

import (
	"atlas-wrg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	ChannelServiceEvent = "channel_service_event"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHANNEL_SERVICE", ChannelServiceEvent, ChannelServerEventCreator(), HandleChannelServerEvent())

}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "World Registry Service", emptyEventCreator, processor)
}
