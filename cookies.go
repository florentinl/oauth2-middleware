package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

func clearStateCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func clearUserCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "oauth_user",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (config Oauth2Config) makeCookie(user string) *http.Cookie {
	expires := time.Now().Local().Add(time.Hour * 1)
	hash := hmac.New(sha256.New, []byte(config.Secret))
	hash.Write([]byte(user))
	hash.Write([]byte(expires.Format(time.UnixDate)))
	signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return &http.Cookie{
		Name:     "oauth_user",
		Value:    user + "|" + expires.Format(time.UnixDate) + "|" + signature,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}
}

func extractCookieInfos(cookie *http.Cookie) (string, string, string, error) {
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 3 {
		return "", "", "", errors.New("Invalid cookie format")
	}
	user := parts[0]
	expires := parts[1]
	signature := parts[2]
	return user, expires, signature, nil
}

func checkExpiration(expires string) error {
	expirationTime, err := time.Parse(time.UnixDate, expires)
	if err != nil {
		return err
	}
	if expirationTime.Before(time.Now()) {
		return errors.New("Cookie expired")
	}
	return nil
}

func checkSignature(user string, expires string, signature string, secret string) bool {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(user))
	hash.Write([]byte(expires))
	expectedSignature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return signature == expectedSignature
}

func (config Oauth2Config) validateCookie(cookie *http.Cookie) (string, error) {
	user, expires, signature, err := extractCookieInfos(cookie)
	if err != nil {
		return "", err
	}

	if err := checkExpiration(expires); err != nil {
		return "", err
	}

	if !checkSignature(user, expires, signature, config.Secret) {
		return "", errors.New("Invalid cookie signature")
	}

	return user, nil
}
