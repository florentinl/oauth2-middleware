package main

import (
	"errors"
	"net/http"
)

func extractCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("oauth_user")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/_auth/login", 302)
		return nil, err
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return cookie, nil
}

func (config Oauth2Config) validate(w http.ResponseWriter, r *http.Request) {
	cookie, err := extractCookie(w, r)
	if err != nil {
		return
	}

	userInfos, err := config.validateCookie(cookie)
	if err != nil {
		clearUserCookie(w)
		http.Redirect(w, r, "/_auth/login", 302)
		return
	}

	w.Header().Set("X-Forwarded-User", userInfos)
	w.WriteHeader(http.StatusOK)
}
