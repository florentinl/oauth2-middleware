package routes

import (
	"log"
	"net/http"

	. "gitlab.viarezo.fr/ViaRezo/oauth2-middleware/utils"
)

func Health(config OAuth2Config, w http.ResponseWriter, r *http.Request) {
	err := config.RedisClient.Ping(config.RedisContext).Err()
	if err != nil {
		log.Println(err)
		http.Error(w, "Lost Connection to Redis", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
