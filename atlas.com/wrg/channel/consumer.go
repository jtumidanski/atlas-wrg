package channel

import (
	"atlas-wrg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus = "channel_service_event"
	topicTokenStatus   = "TOPIC_CHANNEL_SERVICE"
)

func StatusConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[channelServerEvent](consumerNameStatus, topicTokenStatus, groupId, handleStatus())
}

type channelServerEvent struct {
	Status    string `json:"status"`
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	IpAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}

func handleStatus() kafka.HandlerFunc[channelServerEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event channelServerEvent) {
		if event.Status == "STARTED" {
			l.Debugf("Registering channel %d for world %d at %s:%d.", event.ChannelId, event.WorldId, event.IpAddress, event.Port)
			GetChannelRegistry().Register(event.WorldId, event.ChannelId, event.IpAddress, event.Port)
		} else if event.Status == "SHUTDOWN" {
			GetChannelRegistry().RemoveByWorldAndChannel(event.WorldId, event.ChannelId)
		} else {
			l.Errorf("Unhandled event status ", event.Status)
		}
	}
}
