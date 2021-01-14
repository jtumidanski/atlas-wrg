package processor

import (
	"atlas-wrg/attributes"
	"fmt"
	"log"
	"net/http"
)

type Topic struct {
	l *log.Logger
}

func NewTopic(l *log.Logger) *Topic {
	return &Topic{l}
}

func (t *Topic) GetTopic(topic string) (*attributes.TopicData, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/tds/topics/%s", topic))
	if err != nil {
		t.l.Printf("[ERROR] retrieving topic data for %s", topic)
		return nil, err
	}

	td := &attributes.TopicDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		t.l.Printf("[ERROR] decoding topic data for %s", topic)
		return nil, err
	}
	return &td.Data, nil
}
