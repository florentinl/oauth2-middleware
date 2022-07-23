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

// OAuth2 Parameters
type OAuth2Config struct {
	GrantType        string
	ResponseType     string
	Scope            string
	AuthTokenUri     string
	AuthAuthorizeUri string
	AuthUserInfoUri  string
	LogoutUri        string
	OAuth2Clients    map[string]*OAuth2Client
	RedisClient      *redis.Client
	RedisContext     context.Context
}

func getClients() (map[string]*OAuth2Client, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	configFragments := []ConfigFragment{}
	err = json.Unmarshal(file, &configFragments)
	if err != nil {
		return nil, err
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
	return oauth2Clients, nil
}

func NewConfig() OAuth2Config {
	ctx := context.Background()
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       "master",
		SentinelAddrs:    []string{os.Getenv("REDIS_HOST") + ":26379"},
		Password:         os.Getenv("REDIS_PASSWORD"),
		SentinelPassword: os.Getenv("REDIS_PASSWORD"),
		DB:               0,
	})

	oauth2Clients, err := getClients()
	if err != nil {
		log.Fatal(err)
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
		RedisClient:      client,
		RedisContext:     ctx,
	}

	return config
}
