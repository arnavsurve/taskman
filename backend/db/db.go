package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error %s", err)
	}

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("password")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("DB connection successful")

	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountsTable()
}

func (s *PostgresStore) CreateAccountsTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
        id serial primary key,
        username varchar(50),
        password varchar,
        email varchar unique,
        created_at timestamp
    )`

	_, err := s.DB.Exec(query)
	return err
}

// CreateTasksTable creates a new table for a user's tasks with the naming convention t_{id}_{tableName}
func (s *PostgresStore) CreateTasksTable(id, workspaceName string) (string, error) {
	name := fmt.Sprintf("t_%s_%s", id, workspaceName)
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        task_id serial primary key,
        name varchar(50),
		description varchar(255),
		due_date timestamp,
		completion varchar(20) check (completion in ('todo', 'in_progress', 'done')),
		account_id int references accounts(id)
    )`, pq.QuoteIdentifier(name)) // Avoid SQL injection

	_, err := s.DB.Exec(query)
	if err != nil {
		return "", err
	} else {
		fmt.Println("Created:", name)
	}
	return name, err
}

// CreateTask creates a new task in the table with the given name
func (s *PostgresStore) CreateTask(tableName, name, description string, dueDate time.Time, completion shared.CompletionStatus, accountId int) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s(
        name,
        description, 
        due_date, 
        completion, 
        account_id)
        VALUES ($1, $2, $3, $4, $5)`, pq.QuoteIdentifier(tableName))

	_, err := s.DB.Exec(query, name, description, dueDate, completion, accountId)
	if err != nil {
		return "", err
	}

	fmt.Println("Created:", tableName)
	return tableName, nil
}

// GetTasks queries the database and returns a slice of tasks
func (s *PostgresStore) GetTasks(id, tableName string) ([]shared.Task, error) {
	query := fmt.Sprintf(`SELECT task_id, name, description, due_date, completion FROM %s`, tableName)
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []shared.Task

	for rows.Next() {
		task := shared.Task{}
		err := rows.Scan(&task.TaskID, &task.Name, &task.Description, &task.DueDate, &task.CompletionStatus)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// TableExists returns a boolean based on the existence of a table in the database
func (s *PostgresStore) TableExists(tableName string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public'
            AND table_name = $1
        );`

	var exists bool
	err := s.DB.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}
