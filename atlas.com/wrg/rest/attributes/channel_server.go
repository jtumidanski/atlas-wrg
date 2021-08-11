package attributes

type InputChannelServer struct {
	Data ChannelServerData `json:"data"`
}

type ChannelServerListDataContainer struct {
	// All channel servers
	// in: body
	Data []ChannelServerData `json:"data"`
}

type ChannelServerDataContainer struct {
	// A channel server
	// in: body
	Data ChannelServerData `json:"data"`
}

type ChannelServerData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes ChannelServerAttributes `json:"attributes"`
}

type ChannelServerAttributes struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	Capacity  int    `json:"capacity"`
	IpAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}
