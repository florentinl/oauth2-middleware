package internal

import (
	"net/http"
	"net/url"
	"time"
)

func (config OAuth2Config) Login(w http.ResponseWriter, r *http.Request) {
	state, err := RandString(24)
	if err != nil {
		internalServerError(w, err)
		return
	}

	cookie, err := config.makeSession("_auth_state", state, 5*time.Minute)
	if err != nil {
		internalServerError(w, err)
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
