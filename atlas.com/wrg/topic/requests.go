package topic

import (
	"atlas-wrg/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = requests.BaseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

func GetTopic(l logrus.FieldLogger, span opentracing.Span) func(topic string) (*dataBody, error) {
	return func(topic string) (*dataBody, error) {
		td := &dataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(topicById, topic), td, requests.SetRetries(10))
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve topic data for %s", topic)
			return nil, err
		}
		return &td.Data, nil
	}
}