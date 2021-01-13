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
	"atlas-wrg2/attributes"
	"atlas-wrg2/models"
	"atlas-wrg2/registries"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// swagger:route DELETE /channelServers/{id} channelServers unregister
// Removes channel server registration from the world registry
// responses:
//	204: noContentResponse

// GetChannelServers handles DELETE requests
func UnregisterChannelServer(w http.ResponseWriter, r *http.Request) {
	uniqueId := readInt(r, "id")
	registries.GetChannelRegistry().Remove(uniqueId)
	w.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /channelServers channelServers register
// Registers a channel server to the world registry
// responses:
//	200: channelServerResponse

// GetChannelServers handles POST requests
func RegisterChannelServer(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading RegisterChannelServer input")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var input attributes.InputChannelServer
	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		log.Println("Error parsing RegisterChannelServer input")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server := registries.GetChannelRegistry().Register(input.Data.Attributes.WorldId,
		input.Data.Attributes.ChannelId, input.Data.Attributes.IpAddress, input.Data.Attributes.Port)

	var response attributes.ChannelServerDataContainer
	response.Data = getChannelResponseObject(server)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error writing RegisterChannelServer response")
		w.WriteHeader(http.StatusInternalServerError)
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
func GetChannelServers(w http.ResponseWriter, _ *http.Request) {
	var response attributes.ChannelServerListDataContainer
	response.Data = make([]attributes.ChannelServerData, 0)

	for _, x := range registries.GetChannelRegistry().ChannelServers() {
		var serverData = getChannelResponseObject(x)
		response.Data = append(response.Data, serverData)
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding GetChannelServers response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
