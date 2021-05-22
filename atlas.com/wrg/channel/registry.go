package channel

import (
	"sync"
)

type Registry struct {
	mutex   sync.Mutex
	servers []Model
}

var channelRegistry *Registry
var once sync.Once

var uniqueId = 1000000001

func GetChannelRegistry() *Registry {
	once.Do(func() {
		channelRegistry = &Registry{}
	})
	return channelRegistry
}

func (c *Registry) Register(worldId byte, channelId byte, ipAddress string, port int) Model {
	c.mutex.Lock()

	var found *Model = nil
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

	var newChannelServer = NewModel(uniqueId, worldId, channelId, ipAddress, port)
	c.servers = append(c.servers, newChannelServer)
	c.mutex.Unlock()
	return newChannelServer
}

func existingIds(channelServers []Model) []int {
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

func (c *Registry) ChannelServers() []Model {
	servers := c.servers
	return servers
}

func (c *Registry) ChannelServer(worldId byte, channelId byte) *Model {
	for _, x := range c.ChannelServers() {
		if x.WorldId() == worldId && x.ChannelId() == channelId {
			return &x
		}
	}
	return nil
}

func (c *Registry) Remove(id int) {
	c.mutex.Lock()
	index := indexOf(id, c.servers)
	if index >= 0 && index < len(c.servers) {
		c.servers = remove(c.servers, index)
	}
	c.mutex.Unlock()
}

func (c *Registry) RemoveByWorldAndChannel(worldId byte, channelId byte) {
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

func indexOf(uniqueId int, data []Model) int {
	for k, v := range data {
		if uniqueId == v.UniqueId() {
			return k
		}
	}
	return -1 //not found.
}

func remove(s []Model, i int) []Model {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
