package world

import (
	"atlas-wrg/channel"
	"atlas-wrg/configurations"
	"atlas-wrg/json"
	"atlas-wrg/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	getWorlds  = "get_worlds"
	getWorld   = "get_world"
	getChannel = "get_channel"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	wRouter := router.PathPrefix("/worlds").Subrouter()
	wRouter.HandleFunc("/", registerGetWorlds(l)).Methods(http.MethodGet)
	wRouter.HandleFunc("/{worldId}", registerGetWorld(l)).Methods(http.MethodGet)
	wRouter.HandleFunc("/{worldId}/channels/{channelId}", registerGetChannel(l)).Methods(http.MethodGet)
}

func registerGetChannel(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getChannel, func(span opentracing.Span) http.HandlerFunc {
		return parseWorldId(l, func(worldId byte) http.HandlerFunc {
			return parseChannelId(l, func(channelId byte) http.HandlerFunc {
				return GetChannel(l)(span)(worldId, channelId)
			})
		})
	})
}

func registerGetWorld(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getWorld, func(span opentracing.Span) http.HandlerFunc {
		return parseWorldId(l, func(worldId byte) http.HandlerFunc {
			return handleGetWorld(l)(span)(worldId)
		})
	})
}

func registerGetWorlds(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getWorlds, func(span opentracing.Span) http.HandlerFunc {
		return handleGetWorlds(l)(span)
	})
}

type worldIdHandler func(worldId byte) http.HandlerFunc

func parseWorldId(l logrus.FieldLogger, next worldIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		worldId, err := strconv.Atoi(vars["worldId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(byte(worldId))(w, r)
	}
}

type channelIdHandler func(channelId byte) http.HandlerFunc

func parseChannelId(l logrus.FieldLogger, next channelIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId, err := strconv.Atoi(vars["channelId"])
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(byte(channelId))(w, r)
	}
}

func GetChannel(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte, channelId byte) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte, channelId byte) http.HandlerFunc {
		return func(worldId byte, channelId byte) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				server := channel.GetChannelRegistry().ChannelServer(worldId, channelId)
				if server == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				var response channel.DataContainer
				response.Data = getChannelResponseObject(*server)
				err := json.ToJSON(response, w)
				if err != nil {
					l.WithError(err).Errorf("Writing output")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func getChannelResponseObject(server channel.Model) channel.DataBody {
	return channel.DataBody{
		Id:   strconv.Itoa(server.UniqueId()),
		Type: "com.atlas.wrg.rest.attribute.ChannelServerAttributes",
		Attributes: channel.Attributes{
			WorldId:   server.WorldId(),
			ChannelId: server.ChannelId(),
			Capacity:  0,
			IpAddress: server.IpAddress(),
			Port:      server.Port(),
		},
	}
}

func handleGetWorld(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte) http.HandlerFunc {
		return func(worldId byte) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				rd, err := getWorldResponseObject(l, worldId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				response := &DataContainer{Data: *rd}

				err = json.ToJSON(response, w)
				if err != nil {
					l.WithError(err).Errorf("Writing output")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func getWorldResponseObject(l logrus.FieldLogger, worldId byte) (*DataBody, error) {
	c, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		return nil, err
	}

	wc, err := c.GetWorldConfiguration(worldId)
	if err != nil {

		return nil, err
	}

	return &DataBody{
		Id:   strconv.Itoa(int(worldId)),
		Type: "com.atlas.wrg.rest.attribute.WorldAttributes",
		Attributes: Attributes{
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

func handleGetWorlds(l logrus.FieldLogger) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var response DataListContainer
			response.Data = make([]DataBody, 0)

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

			err := json.ToJSON(response, w)
			if err != nil {
				l.WithError(err).Errorf("Writing output")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
