package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
)

func (config OAuth2Config) checkState(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("_auth_state")
	if err != nil {
		return err
	}

	sessionID := cookie.Value
	expectedState, err := config.RedisClient.Get(config.RedisContext, sessionID).Result()
	if err != nil {
		return err
	}

	if r.FormValue("state") != expectedState {
		return CSRFError
	}
	return nil
}

func (config OAuth2Config) getTokens(code string, host string) (Tokens, error) {
	parameters := url.Values{
		"grant_type":    {config.GrantType},
		"client_id":     {config.OAuth2Clients[host].ClientId},
		"client_secret": {config.OAuth2Clients[host].ClientSecret},
		"code":          {code},
		"redirect_uri":  {"https://" + host + "/callback"},
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

func (config OAuth2Config) getUser(token Tokens) (User, error) {
	req, err := http.NewRequest("GET", config.AuthUserInfoUri, nil)
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

func (config OAuth2Config) Callback(w http.ResponseWriter, r *http.Request) {
	err := config.checkState(w, r)
	switch err {
	case http.ErrNoCookie, redis.Nil:
		clearStateCookie(w)
		redirectLogin(w, r)
		return
	case CSRFError:
		http.Error(w, "CSRF validation failed", http.StatusForbidden)
		return
	case nil:
		break
	default:
		internalServerError(w, err)
		return
	}
	clearStateCookie(w)

	tokens, err := config.getTokens(r.FormValue("code"), r.URL.Host)
	if err != nil {
		internalServerError(w, err)
		return
	}

	user, err := config.getUser(tokens)
	if err != nil {
		internalServerError(w, err)
		return
	}

	userStr, err := json.Marshal(user)
	if err != nil {
		internalServerError(w, err)
		return
	}
	userStr64 := base64.StdEncoding.EncodeToString(userStr)
	cookie, err := config.makeSession("_auth_user", userStr64, 5*time.Minute)
	if err != nil {
		internalServerError(w, err)
		return
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "https://"+r.URL.Host+"/", 302)
}
