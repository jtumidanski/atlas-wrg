// Package classification of Channel Server API
//
// Documentation for Channel Server API
//
// Schemes: http
// BasePath: /ms/wrg/channelServers
// Version: 1.0.0
//
// Consumes:
// -application/json
//
// Produces:
// -application/json
// swagger:meta
package handlers

import (
	"atlas-wrg/attributes"
	"atlas-wrg/models"
	"atlas-wrg/registries"
	"log"
	"net/http"
	"strconv"
)

// KeyChannelServer is a key used for the ChannelServer object in the context
type KeyChannelServer struct{}

// Channel Servers handler for getting and updating channel servers
type ChannelServer struct {
	l *log.Logger
}

func NewChannelServer(l *log.Logger) *ChannelServer {
	return &ChannelServer{l}
}

// swagger:route DELETE /channelServers/{channelId} channelServers unregister
// Removes channel server registration from the world registry
// responses:
//	204: noContentResponse

// GetChannelServers handles DELETE requests
func (c *ChannelServer) UnregisterChannelServer(rw http.ResponseWriter, r *http.Request) {
	uniqueId := readInt(r, "channelId")
	registries.GetChannelRegistry().Remove(uniqueId)
	rw.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /channelServers channelServers register
// Registers a channel server to the world registry
// responses:
//	200: channelServerResponse

// GetChannelServers handles POST requests
func (c *ChannelServer) RegisterChannelServer(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(KeyChannelServer{}).(attributes.InputChannelServer)
	server := registries.GetChannelRegistry().Register(input.Data.Attributes.WorldId,
		input.Data.Attributes.ChannelId, input.Data.Attributes.IpAddress, input.Data.Attributes.Port)

	var response attributes.ChannelServerDataContainer
	response.Data = getChannelResponseObject(server)
	err := attributes.ToJSON(response, rw)
	if err != nil {
		c.l.Println("Error writing RegisterChannelServer response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getChannelResponseObject(server models.ChannelServer) attributes.ChannelServerData {
	return attributes.ChannelServerData{
		Id:   strconv.Itoa(server.UniqueId()),
		Type: "com.atlas.wrg.rest.attribute.ChannelServerAttributes",
		Attributes: attributes.ChannelServerAttributes{
			WorldId:   server.WorldId(),
			ChannelId: server.ChannelId(),
			Capacity:  0,
			IpAddress: server.IpAddress(),
			Port:      server.Port(),
		},
	}
}

// swagger:route GET /channelServers channelServers getChannelServers
// Return a list of channel servers in the world registry
// responses:
//	200: channelServersResponse

// GetChannelServers handles GET requests
func (c *ChannelServer) GetChannelServers(rw http.ResponseWriter, _ *http.Request) {
	var response attributes.ChannelServerListDataContainer
	response.Data = make([]attributes.ChannelServerData, 0)

	for _, x := range registries.GetChannelRegistry().ChannelServers() {
		var serverData = getChannelResponseObject(x)
		response.Data = append(response.Data, serverData)
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		c.l.Println("Error encoding GetChannelServers response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
