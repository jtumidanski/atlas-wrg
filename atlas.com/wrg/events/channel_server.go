package events

type ChannelServerEvent struct {
	Status    string
	WorldId   byte
	ChannelId byte
	IpAddress string
	Port      int
}
