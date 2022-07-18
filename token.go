package main

import (
	"net/http"
)

func (config Oauth2Config) token(w http.ResponseWriter, r *http.Request) {
	cookie, err := extractCookie(w, r)
	if err != nil {
		return
	}

	_, err = config.validateCookie(cookie)
	if err != nil {
		clearUserCookie(w)
		http.Redirect(w, r, "/_auth/login", 302)
		return
	}
	http.Redirect(w, r, "http://localhost:8181"+"?token="+cookie.Value, 302)
}
