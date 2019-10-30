package handlers

import (
	"encoding/json"
	userService "github.com/jopicornell/go-rest-api/pkg/api/users/services"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"net/http"
)

var uService *userService.Service = &userService.Service{
	DB: database.GetDB(),
}

func GetUsersHandler(w http.ResponseWriter, _ *http.Request) {
	users, err := uService.GetUsers()
	if err != nil {
		writeErrorToResponse(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		writeErrorToResponse(w, err)
	}
}

func writeErrorToResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	_ = json.NewEncoder(w).Encode(err)
	log.Println(err)
}
