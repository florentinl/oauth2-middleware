package utils

// Structure of the config file
type ConfigFragment struct {
	Name  string   `json:"name"`
	Hosts []string `json:"hosts"`
}

// OAuth2 Client and Secret pair
type OAuth2Client struct {
	ClientId     string
	ClientSecret string
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
