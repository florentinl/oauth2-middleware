package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func (config Oauth2Config) login(w http.ResponseWriter, r *http.Request) {
	state, err := RandString(24)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie, err := config.makeSession("_auth_state", state, 5*time.Minute)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	parameters := url.Values{
		"response_type": {config.ResponseType},
		"client_id":     {config.ClientId},
		"redirect_uri":  {config.BaseUri + "/_auth/callback"},
		"scope":         {config.Scope},
		"state":         {state},
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, config.AuthAuthorizeUri+"?"+parameters.Encode(), 302)
}
