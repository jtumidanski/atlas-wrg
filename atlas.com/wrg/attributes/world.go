package attributes

// A list of worldServers
// swagger:response worldServersResponse
type WorldListDataContainer struct {
	// All world servers
	Data []WorldData `json:"data"`
}

// A worldServer
// swagger:response worldServerResponse
type WorldDataContainer struct {
	// A world server
	Data WorldData `json:"data"`
}

// swagger:model worldServerData
type WorldData struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes WorldAttributes `json:"attributes"`
}

// swagger:model worldServerAttributes
type WorldAttributes struct {
	Name               string `json:"name"`
	Flag               int    `json:"flag"`
	Message            string `json:"message"`
	EventMessage       string `json:"eventMessage"`
	Recommended        bool   `json:"recommended"`
	RecommendedMessage string `json:"recommendedMessage"`
	CapacityStatus     int    `json:"capacityStatus"`
}
