package main

import (
	"fmt"
	"log"
	"net/http"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/internal"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := NewConfig()

	http.HandleFunc("/_auth/login", config.Login)
	http.HandleFunc("/_auth/callback", config.Callback)
	http.HandleFunc("/_auth/validate", config.Validate)
	http.HandleFunc("/_auth/logout", config.Logout)
	http.HandleFunc("/health", config.Health)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
