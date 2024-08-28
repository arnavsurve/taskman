package db

func (s *PostgresStore) CreateWorkspace(name string, accountId int) (int, error) {
	var workspaceID int

	query := `INSERT INTO workspaces (name, account_id)
        VALUES ($1, $2)
        RETURNING workspace_id`
	err := s.DB.QueryRow(query, name, accountId).Scan(&workspaceID)
	if err != nil {
		return 0, err
	}
	return workspaceID, nil
}
