package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GithubLogin(c *gin.Context) {
	url := AppConfig.GitHubLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GithubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusBadRequest, "States don't match!")
		return
	}

	// Get code from the query
	code := c.Query("code")
	fmt.Println(code)
	if code == "" {
		c.String(http.StatusBadRequest, "No code in request")
		return
	}

	// Exchange code for an OAuth token
	config := GithubConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusBadRequest, "Code-token exchange failed")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusBadRequest, "User data fetch failed")
		return
	}
	defer resp.Body.Close()

	// Read user's GitHub data
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "JSON parsing failed")
		return
	}

	c.String(http.StatusOK, string(userData))
}
