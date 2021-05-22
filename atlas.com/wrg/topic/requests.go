package topic

import (
	"atlas-wrg/json"
	"atlas-wrg/rest/requests"
	"atlas-wrg/retry"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = requests.BaseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

func GetTopic(l logrus.FieldLogger) func(topic string) (*dataBody, error) {
	return func(topic string) (*dataBody, error) {
		var r *http.Response
		get := func(attempt int) (bool, error) {
			var err error
			r, err = http.Get(fmt.Sprintf(topicById, topic))
			if err != nil {
				l.Warningln("Unable to retrieve topic data for %s, will retry.", topic)
				return true, err
			}
			return false, nil
		}

		err := retry.Try(get, 10)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve topic data for %s", topic)
			return nil, err
		}

		td := &dataContainer{}
		err = json.FromJSON(td, r.Body)
		if err != nil {
			l.Errorf("Decoding topic data for %s", topic)
			return nil, err
		}
		return &td.Data, nil
	}
}
