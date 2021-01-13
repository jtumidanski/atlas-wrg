package consumers

import (
	"atlas-wrg/attributes"
	"atlas-wrg/events"
	"atlas-wrg/registries"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Consume(ctx context.Context) {
	resp, err := http.Get("http://atlas-nginx:80/ms/tds/topics/TOPIC_CHANNEL_SERVICE")
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var response attributes.TopicDataContainer
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal("Could not unmarshal event into response class ", body)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   response.Data.Attributes.Name,
		GroupID: "World Registry",
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		var event events.ChannelServerEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			processEvent(event)
		}
	}
}

func processEvent(event events.ChannelServerEvent) {
	if event.Status == "STARTED" {
		registries.GetChannelRegistry().Register(event.WorldId, event.ChannelId, event.IpAddress, event.Port)
	} else if event.Status == "SHUTDOWN" {
		registries.GetChannelRegistry().RemoveByWorldAndChannel(event.WorldId, event.ChannelId)
	} else {
		log.Println("Unhandled event status ", event.Status)
	}
}
