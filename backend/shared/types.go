package shared

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CompletionStatus string

const (
	Todo       CompletionStatus = "todo"
	InProgress CompletionStatus = "in_progress"
	Done       CompletionStatus = "done"
)

type Task struct {
	TaskID           int              `json:"task_id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	DueDate          time.Time        `json:"due_date"`   // ex. "due_date": "2023-10-06T15:04:05Z"
	CompletionStatus CompletionStatus `json:"completion"` // todo, in_progress, done
	WorkspaceID      int              `json:"workspace_id"`
	AccountId        int              `json:"account_id"`
}

type Workspace struct {
	WorkspaceID int    `json:"workspace_id"`
	Name        string `json:"name"`
	AccountId   int    `json:"account_id"`
}

type Table struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
