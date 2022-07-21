package internal

import (
	"os"

	"github.com/go-redis/redis"
)

func NewConfig() OAuth2Config {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	config := OAuth2Config{
		GrantType:        "authorization_code",
		ResponseType:     "code",
		Scope:            "default",
		AuthTokenUri:     "https://auth.viarezo.fr/oauth/token",
		AuthAuthorizeUri: "https://auth.viarezo.fr/oauth/authorize",
		AuthAPIUri:       "https://auth.viarezo.fr/api/user/show/me",
		LogoutUri:        "https://auth.viarezo.fr/logout",
		ClientId:         os.Getenv("CLIENT_ID"),
		ClientSecret:     os.Getenv("CLIENT_SECRET"),
		Secret:           os.Getenv("SECRET"),
		BaseUri:          os.Getenv("BASE_URI"),
		RedisClient:      client,
	}

	return config
}
