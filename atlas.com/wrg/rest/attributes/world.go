package attributes

type WorldListDataContainer struct {
	// All world servers
	Data []WorldData `json:"data"`
}

type WorldDataContainer struct {
	// A world server
	Data WorldData `json:"data"`
}

func NewWorldDataContainer(wd WorldData) *WorldDataContainer {
	return &WorldDataContainer{wd}
}

type WorldData struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes WorldAttributes `json:"attributes"`
}

type WorldAttributes struct {
	Name               string `json:"name"`
	Flag               int    `json:"flag"`
	Message            string `json:"message"`
	EventMessage       string `json:"eventMessage"`
	Recommended        bool   `json:"recommended"`
	RecommendedMessage string `json:"recommendedMessage"`
	CapacityStatus     int    `json:"capacityStatus"`
}
