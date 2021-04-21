package channel

import (
	"atlas-wrg/rest/attributes"
	"atlas-wrg/rest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetChannelServers(l *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var response attributes.ChannelServerListDataContainer
		response.Data = make([]attributes.ChannelServerData, 0)

		for _, x := range GetChannelRegistry().ChannelServers() {
			var serverData = getChannelResponseObject(x)
			response.Data = append(response.Data, serverData)
		}

		err := attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func RegisterChannelServer(l *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cs := &attributes.InputChannelServer{}
		err := attributes.FromJSON(cs, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing channel server")
			rw.WriteHeader(http.StatusBadRequest)
			attributes.ToJSON(&resource.GenericError{Message: err.Error()}, rw)
			return
		}

		server := GetChannelRegistry().Register(cs.Data.Attributes.WorldId,
			cs.Data.Attributes.ChannelId, cs.Data.Attributes.IpAddress, cs.Data.Attributes.Port)

		var response attributes.ChannelServerDataContainer
		response.Data = getChannelResponseObject(server)
		err = attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Writing response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func UnregisterChannelServer(l *logrus.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Parsing param channelId as integer")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		GetChannelRegistry().Remove(value)
		rw.WriteHeader(http.StatusNoContent)
	}
}

func getChannelResponseObject(server ChannelServer) attributes.ChannelServerData {
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
