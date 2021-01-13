package handlers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func readInt(r *http.Request, param string) int {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars[param])
	if err != nil {
		log.Println("Error parsing param ", param, " as integer")
		return 0
	}
	return value
}

func readByte(r *http.Request, param string) byte {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars[param])
	if err != nil {
		log.Println("Error parsing param ", param, " as byte")
		return 0
	}
	return byte(value)
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}