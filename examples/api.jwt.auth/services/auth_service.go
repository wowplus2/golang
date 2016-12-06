package services

import (
	"github.com/wowplus2/golang/examples/api.jwt.auth/services/models"
	"net/http"
	"encoding/json"
	"golang.org/x/oauth2/jwt"
)

func Login(reqUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(reqUser) {
		token, err := authBackend.GenerateToken(reqUser.UUID)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			resp, _ := json.Marshal(parameters.TokenAuthentication{token})
			return http.StatusOK, resp
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(reqUser *models.User) []byte {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(reqUser.UUID)
	if err != nil {
		panic(err)
	}
	resp, err := json.Marshal(paramters.TokenAuthentication{token})
	if err != nil {
		panic(err)
	}

	return resp
}

func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}

	tokenString := req.Header.Get("Authorization")

	return authBackend.Logout(tokenString, tokenRequest)
}
