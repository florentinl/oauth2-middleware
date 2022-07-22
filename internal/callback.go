package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (config OAuth2Config) checkState(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("_auth_state")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/login", 302)
		return err
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	sessionID := cookie.Value
	expectedState, err := config.RedisClient.Get(config.RedisContext, sessionID).Result()
	if err != nil {
		log.Println(err)
		clearStateCookie(w)
		http.Redirect(w, r, "/login", 302)
		return err
	}

	if r.FormValue("state") != expectedState {
		http.Error(w, "CSRF validation failed", http.StatusBadRequest)
		return errors.New("CSRF validation failed")
	}
	clearStateCookie(w)
	return nil
}

func (config OAuth2Config) getTokens(code string, baseUri string) (Tokens, error) {
	parameters := url.Values{
		"grant_type":    {config.GrantType},
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {baseUri + "/callback"},
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
	if err != nil {
		return
	}

	tokens, err := config.getTokens(r.FormValue("code"), "https://"+r.URL.Host)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := config.getUser(tokens)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userStr, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userStr64 := base64.StdEncoding.EncodeToString(userStr)
	cookie, err := config.makeSession("_auth_user", userStr64, 5*time.Minute)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "https://"+r.URL.Host+"/validate", 302)
}
