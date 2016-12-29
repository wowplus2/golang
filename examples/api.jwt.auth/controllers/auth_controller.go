package controllers

import (
	"net/http"
	"github.com/wowplus2/golang/examples/api.jwt.auth/services/models"
	"encoding/json"
	"github.com/wowplus2/golang/examples/api.jwt.auth/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	reqUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqUser)

	respStatus, token := services.Login(reqUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	reqUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(services.RefreshToken(reqUser))
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
