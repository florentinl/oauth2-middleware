package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func (config Oauth2Config) login(w http.ResponseWriter, r *http.Request) {
	sessionID, err := RandString(40)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state, err := RandString(24)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = config.RedisClient.Set(sessionID, state, 5*time.Minute).Err()

	cookie := &http.Cookie{
		Name:     "_auth_state",
		Path:     "/",
		Value:    sessionID,
		MaxAge:   int((5 * time.Minute).Seconds()),
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
