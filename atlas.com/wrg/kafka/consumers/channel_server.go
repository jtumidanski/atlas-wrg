package consumers

import (
	"atlas-wrg/channel"
	"github.com/sirupsen/logrus"
)

type channelServerEvent struct {
	Status    string
	WorldId   byte
	ChannelId byte
	IpAddress string
	Port      int
}

func ChannelServerEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &channelServerEvent{}
	}
}

func HandleChannelServerEvent() EventProcessor {
	return func(l logrus.FieldLogger, e interface{}) {
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