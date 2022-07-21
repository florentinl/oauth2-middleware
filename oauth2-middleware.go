package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	config := Oauth2Config{
		GrantType:        "authorization_code",
		ResponseType:     "code",
		Scope:            "default",
		ClientId:         os.Getenv("CLIENT_ID"),
		ClientSecret:     os.Getenv("CLIENT_SECRET"),
		Secret:           os.Getenv("SECRET"),
		BaseUri:          os.Getenv("BASE_URI"),
		AuthTokenUri:     "https://auth.viarezo.fr/oauth/token",
		AuthAuthorizeUri: "https://auth.viarezo.fr/oauth/authorize",
		AuthAPIUri:       "https://auth.viarezo.fr/api/user/show/me",
		LogoutUri:        "https://auth.viarezo.fr/logout",
		RedisClient:      client,
	}
	http.HandleFunc("/_auth/login", config.login)
	http.HandleFunc("/_auth/callback", config.callback)
	http.HandleFunc("/_auth/logout", config.logout)
	http.HandleFunc("/_auth/validate", config.validate)
	http.HandleFunc("/health", config.health)
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
