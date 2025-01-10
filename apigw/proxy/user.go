package proxy

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var (
	loginUrl = "http://usersvc:3001/users/login"
)

type UserLoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var bodyReq UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
			http.Error(w, "Invalid body request", http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(bodyReq)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, loginUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		request.Header.Add("Content-Type", "application/json")
		HandleResponse(w, request)
	}
}
