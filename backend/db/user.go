package db

import (
	"github.com/arnavsurve/taskman/backend/shared"
)

func (s *PostgresStore) CreateGitHubAccount(account *shared.Account) error {
	query := `INSERT INTO accounts(username, email, github_id, oauth_token, created_at)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT DO NOTHING`
	_, err := s.DB.Exec(query, account.Username, account.Email, account.GitHubID, account.OAuthToken, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CheckGitHubUserExists(githubID int) (bool, error) {
	query := `SELECT * FROM accounts where github_id=$1`
	rows, err := s.DB.Query(query, githubID)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
