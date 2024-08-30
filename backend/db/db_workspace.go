package db

import (
	"fmt"

	"github.com/arnavsurve/taskman/backend/shared"
)

func (s *PostgresStore) CreateWorkspace(name string, accountId int) (int, error) {
	var workspaceId int

	query := `INSERT INTO workspaces (name, account_id)
        VALUES ($1, $2)
        RETURNING workspace_id`
	err := s.DB.QueryRow(query, name, accountId).Scan(&workspaceId)
	if err != nil {
		return 0, err
	}
	return workspaceId, nil
}

// GetTasks takes a user ID and the name of the target table and returns the workspace name and a slice of Task structs.
func (s *PostgresStore) GetTasksInWorkspace(id, workspaceId string) (string, []shared.Task, error) {
	query := `SELECT task_id, name, description, due_date, completion, account_id
                            FROM tasks WHERE workspace_id=$1
                            ORDER BY due_date`
	rows, err := s.DB.Query(query, workspaceId)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var tasks []shared.Task
	for rows.Next() {
		task := shared.Task{}
		err := rows.Scan(&task.TaskID, &task.Name, &task.Description, &task.DueDate, &task.CompletionStatus, &task.AccountId)
		if err != nil {
			return "", nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return "", nil, err
	}

	var workspaceName string
	err = s.DB.QueryRow(`SELECT name FROM workspaces WHERE workspace_id=$1`, workspaceId).Scan(&workspaceName)
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}

	return workspaceName, tasks, nil
}

// UpdateWorkspaceByID updates fields (name) of a workspace given account ID, workspace ID, and new workspace name.
func (s *PostgresStore) UpdateWorkspaceByID(accountId, workspaceId int, name string) error {
	query := `UPDATE workspaces SET name=$1 WHERE workspace_id=$2 AND account_id=$3`
	_, err := s.DB.Exec(query, name, workspaceId, accountId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteWorkspaceByID(accountId, workspaceId int) error {
	query := `DELETE FROM tasks WHERE account_id=$1 AND workspace_id=$2`
	_, err := s.DB.Exec(query, accountId, workspaceId)
	if err != nil {
		return err
	}

	query = `DELETE FROM workspaces WHERE account_id=$1 AND workspace_id=$2`
	_, err = s.DB.Exec(query, accountId, workspaceId)
	if err != nil {
		return err
	}

	return nil
}
