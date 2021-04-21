package requests

import (
	"atlas-wrg/rest/attributes"
	"atlas-wrg/retry"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = baseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

var Topic = func(l log.FieldLogger) *topic {
	return &topic{l: l}
}

type topic struct {
	l log.FieldLogger
}

func (t *topic) GetTopic(topic string) (*attributes.TopicData, error) {
	get := func(attempt int) (bool, interface{}, error) {
		r, err := http.Get(fmt.Sprintf(topicById, topic))
		if err != nil {
			t.l.Warningln("Unable to retrieve topic data for %s, will retry.", topic)
			return true, r, err
		}
		return false, r, nil
	}

	r, err := retry.RetryResponse(get, 10)
	if err != nil {
		t.l.Errorf("Unable to retrieve topic data for %s", topic)
		return nil, err
	}
	if val, ok := r.(*http.Response); ok {
		return t.decodeResponse(topic, err, val)
	}
	return nil, errors.New("unexpected output from retry function")
}

func (t *topic) decodeResponse(topic string, err error, val *http.Response) (*attributes.TopicData, error) {
	td := &attributes.TopicDataContainer{}
	err = attributes.FromJSON(td, val.Body)
	if err != nil {
		t.l.Errorf("Decoding topic data for %s", topic)
		return nil, err
	}
	return &td.Data, nil
}
