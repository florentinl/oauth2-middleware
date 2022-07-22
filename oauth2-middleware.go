package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/internal"

	_ "github.com/joho/godotenv/autoload"
)

func substituteXHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Host") != "" {
			r.URL.Host = r.Header.Get("X-Forwarded-Host")
			forwardedUri := strings.Split(r.Header.Get("X-Forwarded-Uri"), "?")
			r.URL.Path = forwardedUri[0]
			if len(forwardedUri) > 1 {
				r.URL.RawQuery = forwardedUri[1]
			}
			r.Form = r.URL.Query()
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	config := NewConfig()

	http.HandleFunc("/", config.Validate)
	http.HandleFunc("/login", config.Login)
	http.HandleFunc("/callback", config.Callback)
	http.HandleFunc("/_logout", config.Logout)
	http.HandleFunc("/health", config.Health)

	xSubstMux := substituteXHeaders(http.DefaultServeMux)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", xSubstMux))
}
