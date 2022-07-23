package internal

import (
	"net/http"
)

func getSessionID(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("_auth_user")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (config OAuth2Config) Validate(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionID(w, r)
	switch err {
	case nil:
		break
	case http.ErrNoCookie:
		redirectLogin(w, r)
		return
	default:
		internalServerError(w, err)
	}

	user, err := config.RedisClient.Get(config.RedisContext, sessionID).Result()
	if err != nil {
		clearUserCookie(w)
		redirectLogin(w, r)
		return
	}

	w.Header().Set("X-Forwarded-User", user)
	w.WriteHeader(http.StatusOK)
}
