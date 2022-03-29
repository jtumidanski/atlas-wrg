package topic

import (
	"atlas-wrg/rest/requests"
	"fmt"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = requests.BaseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

func getTopic(topic string) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(topicById, topic), requests.SetRetries(10))
}
