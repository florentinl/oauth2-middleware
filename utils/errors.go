package utils

import (
	"errors"
	"log"
	"net/http"
)

var CSRFError = errors.New("CSRF validation failed")

func InternalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func RedirectLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.URL.Host+"/_login", 302)
}
