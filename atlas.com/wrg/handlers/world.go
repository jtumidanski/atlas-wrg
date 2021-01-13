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
	"atlas-wrg/attributes"
	"atlas-wrg/configurations"
	"atlas-wrg/models"
	"atlas-wrg/registries"
	"log"
	"net/http"
	"strconv"
)

type World struct {
	l *log.Logger
}

func NewWorld(l *log.Logger) *World {
	return &World{l}
}

// swagger:route GET /worlds/{worldId}/channels/{channelId} worlds getChannelServer
// Retrieves channel server information for a worlds channel.
// responses:
//	200: channelServerResponse
//  404: notFoundResponse

// GetChannel handles GET requests
func (w *World) GetChannel(wr http.ResponseWriter, r *http.Request) {
	worldId := readByte(r, "worldId")
	channelId := readByte(r, "channelId")

	server := registries.GetChannelRegistry().ChannelServer(worldId, channelId)
	if server == nil {
		wr.WriteHeader(http.StatusNotFound)
		return
	}

	var response attributes.ChannelServerDataContainer
	response.Data = getChannelResponseObject(*server)
	err := attributes.ToJSON(response, wr)
	if err != nil {
		log.Println("Error writing GetChannel output")
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /worlds/{worldId} worlds getWorld
// Retrieves world server information.
// responses:
//	200: worldServerResponse
//  404: notFoundResponse

// GetChannel handles GET requests
func (w *World) GetWorld(rw http.ResponseWriter, r *http.Request) {
	worldId := readByte(r, "worldId")
	var response attributes.WorldDataContainer
	response.Data = getWorldResponseObject(worldId)

	err := attributes.ToJSON(response, rw)
	if err != nil {
		log.Println("Error writing GetWorld output")
		rw.WriteHeader(http.StatusInternalServerError)
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
func (w *World) GetWorlds(rw http.ResponseWriter, _ *http.Request) {
	var response attributes.WorldListDataContainer
	response.Data = make([]attributes.WorldData, 0)

	worldIds := mapDistinctWorldId(registries.GetChannelRegistry().ChannelServers())
	for _, id := range worldIds {
		response.Data = append(response.Data, getWorldResponseObject(id))
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		w.l.Println("Error writing GetWorlds output")
		rw.WriteHeader(http.StatusInternalServerError)
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
