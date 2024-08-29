package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

// Init calls each of the table initialization functions and handles errors
func (s *PostgresStore) Init() error {
	initFunctions := []func() error{
		s.CreateAccountsTable,
		s.CreateWorkspacesTable,
		s.CreateTasksTable,
	}

	for _, initFunc := range initFunctions {
		if err := initFunc(); err != nil {
			return err
		}
	}

	return nil
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
	if err != nil {
		return err
	}
	fmt.Println("Successfully created accounts table")
	return nil
}

// CreateWorkspacesTable initializes the table to store information about each workspace
func (s *PostgresStore) CreateWorkspacesTable() error {
	query := `CREATE TABLE IF NOT EXISTS workspaces (
        workspace_id SERIAL PRIMARY KEY,
        name VARCHAR(50),
        account_id INT REFERENCES accounts(id)
        )`

	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created workspaces table")
	return nil
}

// CreateTasksTable creates a new table for a user's tasks with the naming convention t_{id}_{tableName}
func (s *PostgresStore) CreateTasksTable() error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
        task_id serial primary key,
        name varchar(50),
		description varchar(255),
		due_date timestamp,
		completion varchar(20) check (completion in ('todo', 'in_progress', 'done')),
        workspace_id INT references workspaces(workspace_id),
		account_id int references accounts(id)
    )`

	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	} else {
		fmt.Println("Successfully created tasks table")
	}
	return nil
}

func (s *PostgresStore) TaskExists(workspaceId, taskId string) (bool, error) {
	query := `SELECT EXISTS (
        SELECT FROM tasks WHERE workspace_id=$1 AND task_id=$2
    )`

	var exists bool
	err := s.DB.QueryRow(query, workspaceId, taskId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

// WorkspaceExists returns true if a workspace exists in the workspaces table
func (s *PostgresStore) WorkspaceExists(workspaceId string) (bool, error) {
	query := `SELECT EXISTS (
        SELECT FROM workspaces WHERE workspace_id = $1
    )`

	var exists bool
	err := s.DB.QueryRow(query, workspaceId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
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
