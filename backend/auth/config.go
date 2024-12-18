package auth

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Config struct {
	GitHubLoginConfig oauth2.Config
	// add other oauth providers here and configure with a function below
	// then remember to call the function in main (see auth.GithubConfig())
}

var AppConfig Config

func GithubConfig() oauth2.Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %s", err)
	}

	AppConfig.GitHubLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_OAUTH_REDIRECT_URL"),
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}

	return AppConfig.GitHubLoginConfig
}
