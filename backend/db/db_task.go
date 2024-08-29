package db

import (
	"fmt"
	"time"

	"github.com/arnavsurve/taskman/backend/shared"
)

// CreateTask creates a new task in the tasks table
func (s *PostgresStore) CreateTask(name, workspaceId, description string, dueDate time.Time, completion shared.CompletionStatus, accountId int) (string, error) {
	query := `INSERT INTO tasks(
        name,
        description,
        due_date,
        completion,
        workspace_id,
        account_id)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.DB.Exec(query, name, description, dueDate, completion, workspaceId, accountId)
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetTasks takes a user ID and the name of the target table and returns a slice of Task structs.
func (s *PostgresStore) GetTasks(id, workspaceId string) ([]shared.Task, error) {
	query := `SELECT task_id, name, description, due_date, completion, account_id
                            FROM tasks WHERE workspace_id=$1
                            ORDER BY due_date`
	rows, err := s.DB.Query(query, workspaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []shared.Task
	for rows.Next() {
		task := shared.Task{}
		err := rows.Scan(&task.TaskID, &task.Name, &task.Description, &task.DueDate, &task.CompletionStatus, &task.AccountId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID takes a task ID and table name and returns a Task struct
func (s *PostgresStore) GetTaskByID(taskID, workspaceId string) (shared.Task, error) {
	query := `SELECT task_id, name, description, due_date, completion, account_id
						FROM tasks
                        WHERE workspace_id=$1 AND task_id=$2`
	row := s.DB.QueryRow(query, workspaceId, taskID)

	task := shared.Task{}
	err := row.Scan(&task.TaskID, &task.Name, &task.Description, &task.DueDate, &task.CompletionStatus, &task.AccountId)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *PostgresStore) UpdateTaskByID(taskID, workspaceId, name, description string, dueDate time.Time, completion shared.CompletionStatus) error {
	query := `UPDATE tasks set name=$1, description=$2, due_date=$3, completion=$4 where task_id=$5 AND workspace_id=$6`
	_, err := s.DB.Exec(query, name, description, dueDate, completion, taskID, workspaceId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTaskByID takes a task ID and table name and deletes the target row corresponding
// with the target task
func (s *PostgresStore) DeleteTaskByID(taskID, workspaceId string) error {
	query := `DELETE from tasks WHERE task_id=$1 AND workspace_id=$2`
	_, err := s.DB.Exec(query, taskID, workspaceId)
	if err != nil {
		return err
	}
	return nil
}
