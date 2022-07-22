package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Oauth2 Parameters
type OAuth2Config struct {
	GrantType        string
	ResponseType     string
	Scope            string
	AuthTokenUri     string
	AuthAuthorizeUri string
	AuthUserInfoUri  string
	LogoutUri        string
	ClientId         string
	ClientSecret     string
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
