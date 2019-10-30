package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteToJson(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		WriteInternalErrorToResponse(w, err)
	}
}

func WriteInternalErrorToResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	log.Println(err)
}
