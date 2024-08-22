package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

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

func (s *PostgresStore) CreateTasksTable(id string, tableName string) (string, error) {
	name := fmt.Sprintf("t_%s_%s", id, tableName)
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        task_id serial primary key,
		description varchar(255),
		due_date timestamp,
		completion varchar(20) check (completion in ('todo', 'in_progress', 'done')),
		account_id int references accounts(id)
    )`, pq.QuoteIdentifier(name)) // Avoid SQL injection

	_, err := s.DB.Exec(query)
	if err == nil {
		fmt.Println("Created:", name)
	}
	return name, err
}
