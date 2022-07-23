package routes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"
)

func checkState(r *http.Request, redisClient *redis.Client, redisContext context.Context) error {
	cookie, err := r.Cookie("_auth_state")
	if err != nil {
		return err
	}

	sessionID := cookie.Value
	expectedState, err := redisClient.Get(redisContext, sessionID).Result()
	if err != nil {
		return err
	}

	if r.FormValue("state") != expectedState {
		return CSRFError
	}
	return nil
}

func getTokens(config OAuth2Config, code string, host string) (Tokens, error) {
	parameters := url.Values{
		"grant_type":    {config.GrantType},
		"client_id":     {config.OAuth2Clients[host].ClientId},
		"client_secret": {config.OAuth2Clients[host].ClientSecret},
		"code":          {code},
		"redirect_uri":  {"https://" + host + "/_callback"},
		"scope":         {config.Scope},
	}
	resp, err := http.PostForm(config.AuthTokenUri, parameters)
	if err != nil {
		return Tokens{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Tokens{}, errors.New("Invalid response from server")
	}
	defer resp.Body.Close()
	token := Tokens{}
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return Tokens{}, err
	}
	return token, nil
}

func getUser(token Tokens, authUserInfoUri string) (User, error) {
	req, err := http.NewRequest("GET", authUserInfoUri, nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return User{}, errors.New("Identity provider returned an error")
	}
	defer resp.Body.Close()
	user := User{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func Callback(config OAuth2Config, w http.ResponseWriter, r *http.Request) {
	err := checkState(r, config.RedisClient, config.RedisContext)
	switch err {
	case http.ErrNoCookie, redis.Nil:
		RedirectLogin(w, r)
		return
	case CSRFError:
		http.Error(w, "CSRF validation failed", http.StatusForbidden)
		return
	case nil:
		break
	default:
		InternalServerError(w, err)
		return
	}
	ClearStateCookie(w)

	tokens, err := getTokens(config, r.FormValue("code"), r.URL.Host)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	user, err := getUser(tokens, config.AuthUserInfoUri)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	userStr, err := json.Marshal(user)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	userStr64 := base64.StdEncoding.EncodeToString(userStr)
	cookie, err := MakeSession("_auth_user", userStr64, 5*time.Minute, config.RedisClient, config.RedisContext)
	if err != nil {
		InternalServerError(w, err)
		return
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, strings.Join(strings.Split(r.FormValue("state"), ":")[2:], ":"), 302)
}
