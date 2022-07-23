package internal

import (
	"net/http"
	"net/url"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/internal/utils"
)

func (config OAuth2Config) Logout(w http.ResponseWriter, r *http.Request) {
	ClearUserCookie(w)
	redirectLogout := r.FormValue("redirect_logout")
	parameters := url.Values{
		"redirect_logout": {redirectLogout},
	}
	http.Redirect(w, r, config.LogoutUri+"?"+parameters.Encode(), 302)
}
