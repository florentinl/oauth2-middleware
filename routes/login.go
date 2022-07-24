package routes

import (
	"net/http"
	"net/url"
	"time"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"
)

func Login(config OAuth2Config, w http.ResponseWriter, r *http.Request) {
	randString, err := RandString(24)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	state := randString + ":generic-oauth:" + r.FormValue("redirect_login")

	cookie, err := MakeCookie("_auth_state", state, 5*time.Minute, config.RedisClient, config.RedisContext)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	parameters := url.Values{
		"response_type": {config.ResponseType},
		"client_id":     {config.OAuth2Clients[r.URL.Host].ClientId},
		"redirect_uri":  {"https://" + r.URL.Host + "/_callback"},
		"scope":         {config.Scope},
		"state":         {state},
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, config.AuthAuthorizeUri+"?"+parameters.Encode(), 302)
}
