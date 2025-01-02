package board_columns

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BoardColumnsHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterBoardColumnsHandler(router fiber.Router, db *gorm.DB) {
	boardColumnsHandler := BoardColumnsHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/workspace/:workspace_id", boardColumnsHandler.getBoardColumnsByWorkspace)
	router.Get("/:board_column_id/workspace/:workspace_id", boardColumnsHandler.getBoardColumnById)
	router.Post("", boardColumnsHandler.createBoardColumn)
	router.Put("/:board_column_id", boardColumnsHandler.updateBoardColumn)
	router.Delete("/:board_column_id", boardColumnsHandler.deleteBoardColumn)
	//router.Get("/:board_column_id/:field", boardColumnsHandler.getBoardColumnField)
	//router.Put("/:board_column_id/:field", boardColumnsHandler.updateBoardColumnField)
	router.Get("/workspace/:workspace_id/board_column/:board_column_id", boardColumnsHandler.GetSchedulesByBoardColumn)
	router.Get("/:board_column_id", boardColumnsHandler.getBoardColumnById)
	router.Put("/update_position_after_deletion/position", boardColumnsHandler.updatePositionAfterDeletion)
	router.Get("/range/position", boardColumnsHandler.getRage)
	router.Put("/update_position/position", boardColumnsHandler.updatePosition)

}
