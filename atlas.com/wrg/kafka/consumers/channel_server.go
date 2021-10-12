package consumers

import (
	"atlas-wrg/channel"
	"atlas-wrg/kafka/handler"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type channelServerEvent struct {
	Status    string
	WorldId   byte
	ChannelId byte
	IpAddress string
	Port      int
}

func ChannelServerEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &channelServerEvent{}
	}
}

func HandleChannelServerEvent() handler.EventHandler {
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*channelServerEvent); ok {
			if event.Status == "STARTED" {
				channel.GetChannelRegistry().Register(event.WorldId, event.ChannelId, event.IpAddress, event.Port)
			} else if event.Status == "SHUTDOWN" {
				channel.GetChannelRegistry().RemoveByWorldAndChannel(event.WorldId, event.ChannelId)
			} else {
				l.Errorf("Unhandled event status ", event.Status)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
