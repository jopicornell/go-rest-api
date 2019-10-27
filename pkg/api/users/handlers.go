package users

import (
	"encoding/json"
	"log"
	"net/http"
)

func getUsersHandler(w http.ResponseWriter, _ *http.Request) {
	users, err := getUsers()
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		log.Println(err)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}
