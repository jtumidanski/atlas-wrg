package main

import (
	"atlas-wrg/consumers"
	"atlas-wrg/handlers"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	go consumers.Consume(ctx)
	handleRequests()
}

func handleRequests() {
	l := log.New(os.Stdout, "wrg ", log.LstdFlags)

	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/wrg").Subrouter()
	router.Use(commonHeader)
	router.Handle("/docs", middleware.Redoc(middleware.RedocOpts{BasePath: "/ms/wrg", SpecURL: "/ms/wrg/swagger.yaml"}, nil))
	router.Handle("/swagger.yaml", http.StripPrefix("/ms/wrg", http.FileServer(http.Dir("/"))))

	cs := handlers.NewChannelServer(l)
	csRouter := router.PathPrefix("/channelServers").Subrouter()
	csRouter.HandleFunc("/", cs.GetChannelServers).Methods("GET")
	csRouter.Handle("/", cs.MiddlewareValidateChannelServer(cs.RegisterChannelServer)).Methods("POST")
	csRouter.HandleFunc("/{channelId}", cs.UnregisterChannelServer).Methods("DELETE")

	w := handlers.NewWorld(l)
	wRouter := router.PathPrefix("/worlds").Subrouter()
	wRouter.HandleFunc("/", w.GetWorlds).Methods("GET")
	wRouter.HandleFunc("/{worldId}", w.GetWorld).Methods("GET")
	wRouter.HandleFunc("/{worldId}/channels/{channelId}", w.GetChannel).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
