package main

import (
	"log"
	"net/http"
)

func (config Oauth2Config) health(w http.ResponseWriter, r *http.Request) {
	err := config.RedisClient.Ping().Err()
	if err != nil {
		log.Println(err)
		http.Error(w, "Lost Connection to Redis", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
