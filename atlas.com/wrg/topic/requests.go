package topic

import (
	"atlas-wrg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = requests.BaseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

func GetTopic(l logrus.FieldLogger) func(topic string) (*dataBody, error) {
	return func(topic string) (*dataBody, error) {
		td := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(topicById, topic), requests.SetRetries(10))
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve topic data for %s", topic)
			return nil, err
		}
		return &td.Data, nil
	}
}
