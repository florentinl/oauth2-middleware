package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Structure of the config file
type ConfigFile struct {
	Configuration []struct {
		Name  string   `json:"name"`
		Hosts []string `json:"hosts"`
	} `json:"configuration"`
}

// OAuth2 Client and Secret pair
type OAuth2Client struct {
	ClientId     string
	ClientSecret string
}

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
	Secret           string
	RedisClient      *redis.Client
	RedisContext     context.Context
}

// Tokens for a user
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// User is the authenticated user
type User struct {
	ID             int      `json:"id"`
	Login          string   `json:"login"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Email          string   `json:"email"`
	AlternateEmail string   `json:"alternateEmail"`
	BirthDate      string   `json:"birthDate"`
	Promo          int      `json:"promo"`
	Gender         string   `json:"gender"`
	Photo          string   `json:"photo"`
	UpdatedAt      string   `json:"updatedAt"`
	Roles          []string `json:"roles"`
	PersonType     string   `json:"personType"`
}
