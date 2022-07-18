package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

func (config Oauth2Config) checkState(r *http.Request) error {
	state := r.FormValue("state")
	cookie, err := r.Cookie("state")
	if err != nil {
		return err
	}
	stateID := cookie.Value
	state = config.StateMap[stateID]
	if state != r.FormValue("state") || state == "" {
		return errors.New("Mismatching states")
	}
	delete(config.StateMap, stateID)
	return nil
}

func (config Oauth2Config) getTokens(code string) (Tokens, error) {
	// Get the tokens from the user
	parameters := url.Values{
		"grant_type":    {config.GrantType},
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"code":          {code},
		"redirect_uri":  {config.BaseUri + "/_auth/callback"},
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

func (config Oauth2Config) getUser(token Tokens) (User, error) {
	// Get User from API
	req, err := http.NewRequest("GET", config.AuthAPIUri, nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return User{}, errors.New("Invalid response from server")
	}
	defer resp.Body.Close()
	user := User{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (config Oauth2Config) callback(w http.ResponseWriter, r *http.Request) {
	err := config.checkState(r)
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/_auth/login", 302)
		return
	} else if errors.Is(err, errors.New("Mismatching states")) {
		clearStateCookie(w)
		http.Error(w, "Mismatching states", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	clearStateCookie(w)

	tokens, err := config.getTokens(r.FormValue("code"))
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

	// base64 encode the user informations in a signed cookie
	userInfos, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userInfos64 := base64.StdEncoding.EncodeToString(userInfos)
	cookie := config.makeCookie(userInfos64)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/_auth/validate", 302)
}
