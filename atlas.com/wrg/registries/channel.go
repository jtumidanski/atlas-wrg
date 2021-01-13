package registries

import (
	"atlas-wrg2/models"
	"sync"
)

type ChannelRegistry struct {
	mutex   sync.Mutex
	servers []models.ChannelServer
}

var channelRegistry *ChannelRegistry
var once sync.Once

var uniqueId = 1000000001

func GetChannelRegistry() *ChannelRegistry {
	once.Do(func() {
		channelRegistry = &ChannelRegistry{}
	})
	return channelRegistry
}

func (c *ChannelRegistry) Register(worldId byte, channelId byte, ipAddress string, port int) models.ChannelServer {
	c.mutex.Lock()

	var found *models.ChannelServer = nil
	for i := 0; i < len(c.servers); i++ {
		if c.servers[i].WorldId() == worldId && c.servers[i].ChannelId() == channelId {
			found = &c.servers[i]
			break
		}
	}

	if found != nil {
		c.mutex.Unlock()
		return *found
	}

	var existingIds = existingIds(c.servers)

	var currentUniqueId = uniqueId
	for contains(existingIds, currentUniqueId) {
		currentUniqueId = currentUniqueId + 1
		if currentUniqueId > 2000000000 {
			currentUniqueId = 1000000001
		}
		uniqueId = currentUniqueId
	}

	var newChannelServer = models.NewChannelServer(uniqueId, worldId, channelId, ipAddress, port)
	c.servers = append(c.servers, newChannelServer)
	c.mutex.Unlock()
	return newChannelServer
}

func existingIds(channelServers []models.ChannelServer) []int {
	var ids []int
	for _, x := range channelServers {
		ids = append(ids, x.UniqueId())
	}
	return ids
}

func contains(ids []int, id int) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}

func (c *ChannelRegistry) ChannelServers() []models.ChannelServer {
	servers := c.servers
	return servers
}

func (c *ChannelRegistry) ChannelServer(worldId byte, channelId byte) *models.ChannelServer {
	for _, x := range c.ChannelServers() {
		if x.WorldId() == worldId && x.ChannelId() == channelId {
			return &x
		}
	}
	return nil
}

func (c *ChannelRegistry) Remove(id int) {
	c.mutex.Lock()
	index := indexOf(id, c.servers)
	if index >= 0 && index < len(c.servers) {
		c.servers = remove(c.servers, index)
	}
	c.mutex.Unlock()
}

func (c *ChannelRegistry) RemoveByWorldAndChannel(worldId byte, channelId byte) {
	c.mutex.Lock()
	element := c.ChannelServer(worldId, channelId)
	if element != nil {
		index := indexOf(element.UniqueId(), c.servers)
		if index >= 0 && index < len(c.servers) {
			c.servers = remove(c.servers, index)
		}
	}
	c.mutex.Unlock()
}

func indexOf(uniqueId int, data []models.ChannelServer) int {
	for k, v := range data {
		if uniqueId == v.UniqueId() {
			return k
		}
	}
	return -1 //not found.
}

func remove(s []models.ChannelServer, i int) []models.ChannelServer {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
