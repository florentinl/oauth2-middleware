package main

import "github.com/go-redis/redis"

// Oauth2 Parameters
type Oauth2Config struct {
	GrantType        string
	ResponseType     string
	Scope            string
	AuthTokenUri     string
	AuthAuthorizeUri string
	AuthAPIUri       string
	LogoutUri        string
	ClientId         string
	ClientSecret     string
	Secret           string
	BaseUri          string
	RedisClient      *redis.Client
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
