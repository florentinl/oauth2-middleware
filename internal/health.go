package internal

import (
	"log"
	"net/http"
)

func (config OAuth2Config) Health(w http.ResponseWriter, r *http.Request) {
	err := config.RedisClient.Ping(config.RedisContext).Err()
	if err != nil {
		log.Println(err)
		http.Error(w, "Lost Connection to Redis", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
