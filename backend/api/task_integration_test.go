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

// type MockPostgresStore struct {
// 	MockTableExists func(tableName string) (bool, error)
// 	MockGetTasks    func(id, tableName string) ([]shared.Task, error)
// }
//
// func (m *MockPostgresStore) TableExists(tableName string) (bool, error) {
// 	return m.MockTableExists(tableName)
// }
//
// func (m *MockPostgresStore) GetTasks(id, tableName string) (bool, error) {
// 	return m.MockTableExists(tableName)
// }
//
// func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
//
// }

func TestHandleCreateTask(t *testing.T) {

}
