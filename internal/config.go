package internal

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

	// Parse Config File as JSON
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	configFragments := []ConfigFragment{}
	err = json.Unmarshal(file, &configFragments)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	oauth2Clients := make(map[string]*OAuth2Client)
	for _, clientConfig := range configFragments {
		client := OAuth2Client{
			ClientId:     os.Getenv(strings.ToUpper(clientConfig.Name) + "_CLIENT_ID"),
			ClientSecret: os.Getenv(strings.ToUpper(clientConfig.Name) + "_CLIENT_SECRET"),
		}
		for _, host := range clientConfig.Hosts {
			oauth2Clients[host] = &client
		}
	}

	config := OAuth2Config{
		GrantType:        "authorization_code",
		ResponseType:     "code",
		Scope:            "default",
		AuthTokenUri:     os.Getenv("OAUTH_TOKEN_URI"),
		AuthAuthorizeUri: os.Getenv("OAUTH_AUTHORIZE_URI"),
		AuthUserInfoUri:  os.Getenv("OAUTH_USERINFO_URI"),
		LogoutUri:        os.Getenv("OAUTH_LOGOUT_URI"),
		OAuth2Clients:    oauth2Clients,
		Secret:           os.Getenv("SECRET"),
		RedisClient:      client,
		RedisContext:     ctx,
	}

	log.Println("[INFO] OAuth2Config:", config)

	return config
}
