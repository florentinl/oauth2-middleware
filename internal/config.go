package internal

import (
	"context"
	"log"
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
		AuthTokenUri:     os.Getenv("OAUTH_TOKEN_URI"),
		AuthAuthorizeUri: os.Getenv("OAUTH_AUTHORIZE_URI"),
		AuthUserInfoUri:  os.Getenv("OAUTH_USERINFO_URI"),
		LogoutUri:        os.Getenv("OAUTH_LOGOUT_URI"),
		ClientId:         os.Getenv("CLIENT_ID"),
		ClientSecret:     os.Getenv("CLIENT_SECRET"),
		Secret:           os.Getenv("SECRET"),
		RedisClient:      client,
		RedisContext:     ctx,
	}

	log.Println("[INFO] OAuth2Config:", config)

	return config
}
