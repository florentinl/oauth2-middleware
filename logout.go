package main

import (
	"net/http"
	"net/url"
)

func (config Oauth2Config) logout(w http.ResponseWriter, r *http.Request) {
	clearUserCookie(w)
	redirectLogout := r.FormValue("redirect_logout")
	parameters := url.Values{
		"redirect_logout": {redirectLogout},
	}
	http.Redirect(w, r, config.LogoutUri+"?"+parameters.Encode(), 302)
}
