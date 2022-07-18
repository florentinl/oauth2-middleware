package main

import (
	"net/http"
	"net/url"
)

func (config Oauth2Config) logout(w http.ResponseWriter, r *http.Request) {
	clearUserCookie(w)
	parameters := url.Values{
		"redirect_logout": {config.BaseUri + "/"},
	}
	http.Redirect(w, r, config.LogoutUri+"?"+parameters.Encode(), 302)
}
