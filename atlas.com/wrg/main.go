package main

import (
	"atlas-wrg2/consumers"
	"atlas-wrg2/handlers"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func main() {
	ctx := context.Background()
	go consumers.Consume(ctx)
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonHeader)
	router.HandleFunc("/ms/wrg/channelServers", handlers.GetChannelServers).Methods("GET")
	router.HandleFunc("/ms/wrg/channelServers", handlers.RegisterChannelServer).Methods("POST")
	router.HandleFunc("/ms/wrg/channelServers/{id}", handlers.UnregisterChannelServer).Methods("DELETE")
	router.HandleFunc("/ms/wrg/worlds", handlers.GetWorlds).Methods("GET")
	router.HandleFunc("/ms/wrg/worlds/{worldId}", handlers.GetWorld).Methods("GET")
	router.HandleFunc("/ms/wrg/worlds/{worldId}/channels/{channelId}", handlers.GetChannel).Methods("GET")

	ops := middleware.RedocOpts{BasePath: "/ms/wrg", SpecURL: "/ms/wrg/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	router.Handle("/ms/wrg/docs", sh)

	rewrite := func(path string) string {
		// your rewrite code, returns the new path
		return strings.ReplaceAll(path, "/ms/wrg", "")
	}

	fileServer := http.FileServer(http.Dir("/"))
	router.HandleFunc("/ms/wrg/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = rewrite(r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}



func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
