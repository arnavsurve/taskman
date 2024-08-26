package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type MockPostgresStore struct {
	MockTableExists func(tableName string) (bool, error)
	MockGetTasks    func(id, tableName string) ([]shared.Task, error)
}

func (m *MockPostgresStore) TableExists(tableName string) (bool, error) {
	return m.MockTableExists(tableName)
}

func (m *MockPostgresStore) GetTasks(id, tableName string) (bool, error) {
	return m.MockTableExists(tableName)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	return ctx, w
}

func TestHandleGetTasks_Success(t *testing.T) {
	ctx, w := setupTestContext()

	// Mock user claims
	ctx.Set("user", map[string]interface{}{"id": 1.0})

	// Mock URL parameters
	ctx.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "workspace", Value: "test_workspace"},
	}

	// Mock store
	mockStore := &MockPostgresStore{
		MockTableExists: func(tableName string) (bool, error) {
			return true, nil
		},
		MockGetTasks: func(id, tableName string) ([]shared.Task, error) {
			return []shared.Task{
				{TaskID: 1, Name: "Test Task", Description: "Test Description"},
			}, nil
		},
	}

	HandleGetTasks(ctx, mockStore)
}
