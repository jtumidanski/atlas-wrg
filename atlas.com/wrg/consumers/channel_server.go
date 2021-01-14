package consumers

import (
	"atlas-wrg/events"
	"atlas-wrg/processor"
	"atlas-wrg/registries"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type ChannelServer struct {
	l   *log.Logger
	ctx context.Context
}

func NewChannelServer(l *log.Logger, ctx context.Context) *ChannelServer {
	return &ChannelServer{l, ctx}
}

func (c *ChannelServer) Init() {
	t := processor.NewTopic(c.l)
	td, err := t.GetTopic("TOPIC_CHANNEL_SERVICE")
	if err != nil {
		c.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	fmt.Print(td)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: "World Registry",
		MaxWait: 50 * time.Millisecond,
	})
	for {
		msg, err := r.ReadMessage(c.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		var event events.ChannelServerEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			c.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			c.processEvent(event)
		}
	}
}

func (c *ChannelServer) processEvent(event events.ChannelServerEvent) {
	if event.Status == "STARTED" {
		registries.GetChannelRegistry().Register(event.WorldId, event.ChannelId, event.IpAddress, event.Port)
	} else if event.Status == "SHUTDOWN" {
		registries.GetChannelRegistry().RemoveByWorldAndChannel(event.WorldId, event.ChannelId)
	} else {
		c.l.Println("Unhandled event status ", event.Status)
	}
}
