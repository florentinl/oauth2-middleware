package main

import (
	"log"
	"net/http"
	"net/url"
)

func (config Oauth2Config) login(w http.ResponseWriter, r *http.Request) {
	state, err := RandString(24)
	if err != nil {
		log.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:     "state",
		Path:     "/",
		Value:    state,
		MaxAge:   1200,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	clearUserCookie(w)

	parameters := url.Values{
		"response_type": {config.ResponseType},
		"client_id":     {config.ClientId},
		"redirect_uri":  {config.BaseUri + "/_auth/callback"},
		"scope":         {config.Scope},
		"state":         {state},
	}

	http.Redirect(w, r, config.AuthAuthorizeUri+"?"+parameters.Encode(), 302)
}
