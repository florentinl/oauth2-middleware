package internal

import (
	"errors"
	"net/http"
)

func getSessionID(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("_auth_user")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "https://"+getBaseUri(r)+"/_auth/login", 302)
		return "", err
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return "", err
	}
	return cookie.Value, nil
}

func (config OAuth2Config) Validate(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionID(w, r)
	if err != nil {
		return
	}

	user, err := config.RedisClient.Get(sessionID).Result()
	if err != nil {
		clearUserCookie(w)
		http.Redirect(w, r, "https://"+getBaseUri(r)+"/_auth/login", 302)
		return
	}

	w.Header().Set("X-Forwarded-User", user)
	w.WriteHeader(http.StatusOK)
}
