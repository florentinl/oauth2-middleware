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
		MaxAge:   -int(24 * time.Hour.Seconds()),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func ClearUserCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "_auth_user",
		Value:    "",
		Path:     "/",
		MaxAge:   -int(24 * time.Hour.Seconds()),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func DeleteSession(redisClient *redis.Client, redisContext context.Context, r *http.Request) error {
	cookie, err := r.Cookie("_auth_user")
	if err != nil {
		return err
	}
	return redisClient.Del(redisContext, cookie.Value).Err()
}

func MakeCookie(name string, payload string, maxAge time.Duration, redisClient *redis.Client, redisContext context.Context) (*http.Cookie, error) {
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
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}, nil
}

func MakeSession(name string, payload string, redisClient *redis.Client, redisContext context.Context) (*http.Cookie, error) {
	sessionID, err := RandString(40)
	if err != nil {
		return nil, err
	}

	err = redisClient.Set(redisContext, sessionID, payload, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     name,
		Path:     "/",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
	}, nil
}
