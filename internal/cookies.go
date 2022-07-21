package internal

import (
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

func (config Oauth2Config) makeSession(name string, payload string, maxAge time.Duration) (*http.Cookie, error) {
	sessionID, err := RandString(40)
	if err != nil {
		return nil, err
	}

	err = config.RedisClient.Set(sessionID, payload, maxAge).Err()
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     name,
		Path:     "/",
		Value:    sessionID,
		MaxAge:   int((maxAge).Seconds()),
		HttpOnly: true,
		Secure:   true,
	}, nil
}
