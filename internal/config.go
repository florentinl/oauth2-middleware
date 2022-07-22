package internal

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewConfig() OAuth2Config {
	ctx := context.Background()
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       "master",
		SentinelAddrs:    []string{os.Getenv("REDIS_HOST") + ":26379"},
		Password:         os.Getenv("REDIS_PASSWORD"),
		SentinelPassword: os.Getenv("REDIS_PASSWORD"),
		DB:               0,
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
		RedisClient:      client,
		RedisContext:     ctx,
	}

	return config
}
