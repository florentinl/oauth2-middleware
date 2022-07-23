package main

import (
	"fmt"
	"log"
	"net/http"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/routes"
	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := NewConfig()

	http.HandleFunc("/", ValidateHandler(config))
	http.HandleFunc("/_login", LoginHandler(config))
	http.HandleFunc("/_callback", CallbackHandler(config))
	http.HandleFunc("/_logout", LogoutHandler(config))
	http.HandleFunc("/_health", HealthHandler(config))

	xSubstMux := SubstituteXHeaders(http.DefaultServeMux)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", xSubstMux))
}
