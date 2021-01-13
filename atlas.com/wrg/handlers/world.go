// Package classification of World API
//
// Documentation for World API
//
// Schemes: http
// BasePath: /ms/wrg/worlds
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
	"atlas-wrg2/configurations"
	"atlas-wrg2/models"
	"atlas-wrg2/registries"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// swagger:route GET /worlds/{worldId}/channels/{channelId} worlds getChannelServer
// Retrieves channel server information for a worlds channel.
// responses:
//	200: channelServerResponse
//  404: notFoundResponse

// GetChannel handles GET requests
func GetChannel(w http.ResponseWriter, r *http.Request) {
	worldId := readByte(r, "worldId")
	channelId := readByte(r, "channelId")

	server := registries.GetChannelRegistry().ChannelServer(worldId, channelId)
	if server == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var response attributes.ChannelServerDataContainer
	response.Data = getChannelResponseObject(*server)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error writing GetChannel output")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /worlds/{worldId} worlds getWorld
// Retrieves world server information.
// responses:
//	200: worldServerResponse
//  404: notFoundResponse

// GetChannel handles GET requests
func GetWorld(w http.ResponseWriter, r *http.Request) {
	worldId := readByte(r, "worldId")
	var response attributes.WorldDataContainer
	response.Data = getWorldResponseObject(worldId)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error writing GetWorld output")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getWorldResponseObject(worldId byte) attributes.WorldData {
	c, _ := configurations.GetConfiguration()
	wc, _ := c.GetWorldConfiguration(worldId)

	return attributes.WorldData{
		Id:   strconv.Itoa(int(worldId)),
		Type: "com.atlas.wrg.rest.attribute.WorldAttributes",
		Attributes: attributes.WorldAttributes{
			Name:               wc.Name,
			Flag:               getFlag(wc.Flag),
			Message:            wc.ServerMessage,
			EventMessage:       wc.EventMessage,
			Recommended:        wc.WhyAmIRecommended != "",
			RecommendedMessage: wc.WhyAmIRecommended,
			CapacityStatus:     0,
		},
	}
}

// swagger:route GET /worlds worlds getWorlds
// Retrieves all world server information.
// responses:
//	200: worldServersResponse

// GetChannel handles GET requests
func GetWorlds(w http.ResponseWriter, _ *http.Request) {
	var response attributes.WorldListDataContainer
	response.Data = make([]attributes.WorldData, 0)

	worldIds := mapDistinctWorldId(registries.GetChannelRegistry().ChannelServers())
	for _, id := range worldIds {
		response.Data = append(response.Data, getWorldResponseObject(id))
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error writing GetWorlds output")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func mapDistinctWorldId(channelServers []models.ChannelServer) []byte {
	m := make(map[byte]struct{})
	for _, element := range channelServers {
		m[element.WorldId()] = struct{}{}
	}

	keys := make([]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getFlag(flag string) int {
	switch flag {
	case "NOTHING":
		return 0
	case "EVENT":
		return 1
	case "NEW":
		return 2
	case "HOT":
		return 3
	default:
		return 0
	}
}
