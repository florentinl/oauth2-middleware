package routes

import (
	"net/http"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"
)

func LoginHandler(config OAuth2Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Login(config, w, r)
	}
}

func CallbackHandler(config OAuth2Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Callback(config, w, r)
	}
}

func ValidateHandler(config OAuth2Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Validate(config, w, r)
	}
}

func LogoutHandler(config OAuth2Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Logout(config, w, r)
	}
}

func HealthHandler(config OAuth2Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Health(config, w, r)
	}
}
