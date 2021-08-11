package world

import (
	"atlas-wrg/channel"
	"atlas-wrg/configurations"
	"atlas-wrg/rest/attributes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetChannel(l logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId := byte(value)

		vars = mux.Vars(r)
		value, err = strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		channelId := byte(value)

		server := channel.GetChannelRegistry().ChannelServer(worldId, channelId)
		if server == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var response attributes.ChannelServerDataContainer
		response.Data = getChannelResponseObject(*server)
		err = attributes.ToJSON(response, w)
		if err != nil {
			l.WithError(err).Errorf("Writing GetChannel output")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func getChannelResponseObject(server channel.Model) attributes.ChannelServerData {
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

func GetWorld(l logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		worldId := byte(value)

		rd, err := getWorldResponseObject(l, worldId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response := attributes.NewWorldDataContainer(*rd)

		err = attributes.ToJSON(response, w)
		if err != nil {
			l.WithError(err).Errorf("Writing GetWorld output")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getWorldResponseObject(l logrus.FieldLogger, worldId byte) (*attributes.WorldData, error) {
	c, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		return nil, err
	}

	wc, err := c.GetWorldConfiguration(worldId)
	if err != nil {

		return nil, err
	}

	return &attributes.WorldData{
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
	}, nil
}

func GetWorlds(l logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response attributes.WorldListDataContainer
		response.Data = make([]attributes.WorldData, 0)

		worldIds := mapDistinctWorldId(channel.GetChannelRegistry().ChannelServers())
		for _, id := range worldIds {
			rd, err := getWorldResponseObject(l, id)
			if err != nil {
				l.WithError(err).Errorf("Unable to get response object")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.Data = append(response.Data, *rd)
		}

		err := attributes.ToJSON(response, w)
		if err != nil {
			l.WithError(err).Errorf("Writing GetWorlds output")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func mapDistinctWorldId(channelServers []channel.Model) []byte {
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
