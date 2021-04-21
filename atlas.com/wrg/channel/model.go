package channel

type ChannelServer struct {
	uniqueId  int
	worldId   byte
	channelId byte
	ipAddress string
	port      int
}

func (c ChannelServer) UniqueId() int {
	return c.uniqueId
}

func (c ChannelServer) WorldId() byte {
	return c.worldId
}

func (c ChannelServer) ChannelId() byte {
	return c.channelId
}

func (c ChannelServer) IpAddress() string {
	return c.ipAddress
}

func (c ChannelServer) Port() int {
	return c.port
}

func NewChannelServer(uniqueId int, worldId byte, channelId byte, ipAddress string, port int) ChannelServer {
	return ChannelServer{
		uniqueId:  uniqueId,
		worldId:   worldId,
		channelId: channelId,
		ipAddress: ipAddress,
		port:      port,
	}
}
