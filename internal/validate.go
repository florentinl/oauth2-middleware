package internal

import (
	"net/http"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/internal/utils"
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
		RedirectLogin(w, r)
		return
	default:
		InternalServerError(w, err)
	}

	user, err := config.RedisClient.Get(config.RedisContext, sessionID).Result()
	if err != nil {
		ClearUserCookie(w)
		RedirectLogin(w, r)
		return
	}

	w.Header().Set("X-Forwarded-User", user)
	w.WriteHeader(http.StatusOK)
}
