package routes

import (
	"net/http"
	"net/url"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"
)

func Logout(config OAuth2Config, w http.ResponseWriter, r *http.Request) {
	ClearUserCookie(w)
	redirectLogout := r.FormValue("redirect_logout")
	parameters := url.Values{
		"redirect_logout": {redirectLogout},
	}
	http.Redirect(w, r, config.LogoutUri+"?"+parameters.Encode(), 302)
}
