package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

func ClearStateCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "_auth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func ClearUserCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "_auth_user",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func MakeSession(name string, payload string, maxAge time.Duration, redisClient *redis.Client, redisContext context.Context) (*http.Cookie, error) {
	sessionID, err := RandString(40)
	if err != nil {
		return nil, err
	}

	err = redisClient.Set(redisContext, sessionID, payload, maxAge).Err()
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
