package utils

import (
	"log"
	"net/http"
)

func AuthError(w http.ResponseWriter, err error) {
	log.Printf(err.Error())
	w.WriteHeader(401)
	w.Write([]byte(err.Error()))
}
