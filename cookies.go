package main

import (
	"log"
	"net/http"
	"time"
)

func clearStateCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "_auth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func clearUserCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "_auth_user",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (config Oauth2Config) makeCookie(user string) (*http.Cookie, error) {
	sessionID, err := RandString(40)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = config.RedisClient.Set(sessionID, user, 5*time.Minute).Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &http.Cookie{
		Name:     "_auth_user",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int((5 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   true,
	}, nil
}
