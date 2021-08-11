package channel

import (
	"atlas-wrg/json"
	"atlas-wrg/rest/resource"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetChannelServers(l logrus.FieldLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var response DataListContainer
		response.Data = make([]DataBody, 0)

		for _, x := range GetChannelRegistry().ChannelServers() {
			var serverData = getChannelResponseObject(x)
			response.Data = append(response.Data, serverData)
		}

		err := json.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Encoding response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func RegisterChannelServer(l logrus.FieldLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cs := &InputDataContainer{}
		err := json.FromJSON(cs, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing channel server")
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&resource.GenericError{Message: err.Error()}, rw)
			return
		}

		server := GetChannelRegistry().Register(cs.Data.Attributes.WorldId,
			cs.Data.Attributes.ChannelId, cs.Data.Attributes.IpAddress, cs.Data.Attributes.Port)

		var response DataContainer
		response.Data = getChannelResponseObject(server)
		err = json.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Writing response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func UnregisterChannelServer(l logrus.FieldLogger) http.HandlerFunc {
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

func getChannelResponseObject(server Model) DataBody {
	return DataBody{
		Id:   strconv.Itoa(server.UniqueId()),
		Type: "com.atlas.wrg.rest.attribute.ChannelServerAttributes",
		Attributes: Attributes{
			WorldId:   server.WorldId(),
			ChannelId: server.ChannelId(),
			Capacity:  0,
			IpAddress: server.IpAddress(),
			Port:      server.Port(),
		},
	}
}
