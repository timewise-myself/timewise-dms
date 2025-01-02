package workspace_log

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWorkspaceLogHandler(router fiber.Router, db *gorm.DB) {
	workspaceLogHandler := WorkspaceLog{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", workspaceLogHandler.getWorkspaceLog)
	router.Get("/:workspace_log_id", workspaceLogHandler.getWorkspaceLogById)
	router.Delete("/:workspace_log_id", workspaceLogHandler.removeWorkspaceLogById)
	router.Post("/", workspaceLogHandler.createWorkspaceLog)
	router.Get("/workspace/:workspace_id", workspaceLogHandler.getWorkspaceLogsByWorkspaceId)
}
