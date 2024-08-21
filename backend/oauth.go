package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %s", err)
	}
}

// oauthConfig represents the OAuth2 configuration for GitHub authentication.
var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	Scopes:       []string{"read:user", "user:email"},
	Endpoint:     github.Endpoint,
}
