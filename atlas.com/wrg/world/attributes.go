package world

type DataListContainer struct {
	// All world servers
	Data []DataBody `json:"data"`
}

type DataContainer struct {
	// A world server
	Data DataBody `json:"data"`
}

type DataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name               string `json:"name"`
	Flag               int    `json:"flag"`
	Message            string `json:"message"`
	EventMessage       string `json:"eventMessage"`
	Recommended        bool   `json:"recommended"`
	RecommendedMessage string `json:"recommendedMessage"`
	CapacityStatus     int    `json:"capacityStatus"`
}
