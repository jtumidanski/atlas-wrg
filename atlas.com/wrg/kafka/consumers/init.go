package consumers

import (
	"atlas-wrg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHANNEL_SERVICE", ChannelServerEventCreator(), HandleChannelServerEvent())

}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, topicToken, "World Registry Service", emptyEventCreator, processor)
}
