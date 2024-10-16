package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnavsurve/taskman/backend/db"
	"github.com/arnavsurve/taskman/backend/shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func HandleCreateWorkspace(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Binding JSON body to a Workspace struct
	workspace := shared.Workspace{AccountId: userID}
	if err := ctx.ShouldBindJSON(&workspace); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workspace.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace name cannot be empty"})
		return
	}
	if workspace.AccountId < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	workspaceID, err := store.CreateWorkspace(workspace.Name, workspace.AccountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Successfully created workspace", "workspace_id": workspaceID})
}

// HandleGetTasksFromWorkspace takes url parameters ID and workspace name. It calls GetTasks and returns a
// JSON object holding task JSON objects.
func HandleGetTasksInWorkspace(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// Checking if the requested workspace for the task exists
	workspaceId := ctx.Param("workspaceId")
	exists, err := store.WorkspaceExists(workspaceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if workspace exists"})
		fmt.Println(err)
		return
	}
	if exists != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace does not exist"})
		return
	}

	workspaceName, tasks, err := store.GetTasksInWorkspace(requestedID, workspaceId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{workspaceName: tasks})
}

func HandleUpdateWorkspaceByID(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	workspaceId := ctx.Param("workspaceId")
	intWorkspaceId, _ := strconv.Atoi(workspaceId)

	// Binding JSON body to a Workspace struct
	workspace := shared.Workspace{AccountId: userID, WorkspaceID: intWorkspaceId}
	if err := ctx.ShouldBindJSON(&workspace); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workspace.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Workspace name cannot be empty"})
		return
	}
	if workspace.AccountId < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = store.UpdateWorkspaceByID(workspace.AccountId, workspace.WorkspaceID, workspace.Name)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Workspace name successfully updated"})
}

func HandleDeleteWorkspaceByID(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	workspaceId := ctx.Param("workspaceId")
	intWorkspaceId, _ := strconv.Atoi(workspaceId)

	err = store.DeleteWorkspaceByID(intRequestedID, intWorkspaceId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Workspace successfully deleted"})
}

func HandleListWorkspaces(ctx *gin.Context, store *db.PostgresStore) {
	userClaims := ctx.MustGet("user").(jwt.MapClaims)
	userID := int(userClaims["id"].(float64))

	// Verify user's ID matches ID of the resource
	requestedID := ctx.Param("id")
	intRequestedID, err := strconv.Atoi(requestedID)
	if err != nil || intRequestedID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	workspaces, err := store.ListWorkspaces(requestedID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"workspaces": workspaces})
}
