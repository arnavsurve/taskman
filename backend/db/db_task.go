package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// CreateTask creates a new task in the tasks table
func (s *PostgresStore) CreateTask(name, description string, dueDate time.Time, completion shared.CompletionStatus, workspaceId, accountId int) (string, error) {
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

	fmt.Println("Created task", name, "in workspace", strconv.Itoa(workspaceId))
	return name, nil
}

// GetTasks takes a user ID and the name of the target table and returns a slice of Task structs.
func (s *PostgresStore) GetTasks(id, tableName string) ([]shared.Task, error) {
	query := fmt.Sprintf(`SELECT task_id, name, description, due_date, completion, account_id
                            FROM %s
                            ORDER BY due_date`, pq.QuoteIdentifier(tableName))
	rows, err := s.DB.Query(query)
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
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID takes a task ID and table name and returns a Task struct
func (s *PostgresStore) GetTaskByID(taskID, tableName string) (shared.Task, error) {
	query := fmt.Sprintf(`SELECT task_id, name, description, due_date, completion, account_id
                            FROM %s WHERE task_id = %s`, pq.QuoteIdentifier(tableName), taskID)
	row := s.DB.QueryRow(query)

	task := shared.Task{}
	err := row.Scan(&task.TaskID, &task.Name, &task.Description, &task.DueDate, &task.CompletionStatus, &task.AccountId)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (s *PostgresStore) UpdateTaskByID(taskID, tableName, name, description string, dueDate time.Time, completion shared.CompletionStatus) error {
	query := fmt.Sprintf(`UPDATE %s set name=$1, description=$2, due_date=$3, completion=$4 where task_id=$5`, pq.QuoteIdentifier(tableName))
	_, err := s.DB.Exec(query, name, description, dueDate, completion, taskID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTaskByID takes a task ID and table name and deletes the target row corresponding
// with the target task
func (s *PostgresStore) DeleteTaskByID(taskID, tableName string) error {
	query := fmt.Sprintf(`DELETE from %s WHERE task_id = %s`, pq.QuoteIdentifier(tableName), taskID)
	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
