package attributes

type TopicDataContainer struct {
	Data TopicData `json:"data"`
}

type TopicData struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes TopicAttributes `json:"attributes"`
}

type TopicAttributes struct {
	Name               string `json:"name"`
}
