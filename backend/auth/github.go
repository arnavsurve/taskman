package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/arnavsurve/taskman/backend/utils"
	"github.com/gin-gonic/gin"
)

func GithubLogin(c *gin.Context) {
	url := AppConfig.GitHubLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GithubCallback(c *gin.Context, store *db.PostgresStore) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusBadRequest, "States don't match!")
		return
	}

	// Get code from the query
	code := c.Query("code")
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

	// Variable to hold the parsed data
	var githubAccount shared.GitHubAccount
	// Parse (unmarshal) the JSON into the struct
	err = json.Unmarshal([]byte(userData), &githubAccount)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	githubAccount.OAuthToken = token.AccessToken

	// Map GitHubAccount fields to Account struct
	account := shared.NewGitHubAccount(githubAccount)

	// Check if user already exists in accounts table
	if exists, _ := store.CheckGitHubUserExists(githubAccount.GitHubID); exists == false {
		// Create entry in accounts table for this user
		err = store.CreateGitHubAccount(account)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Retrieve user ID
	var userID int
	query := `select id from accounts where github_id= $1`
	err = store.DB.QueryRow(query, githubAccount.GitHubID).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	// Generate JWT token
	JWT, err := utils.GenerateToken(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("JWT: %s\n", JWT)
	c.JSON(http.StatusOK, gin.H{"token": JWT})
}
