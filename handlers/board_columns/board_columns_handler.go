package board_columns

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

// getBoardColumnsByWorkspace godoc
// @Summary Get board columns by workspace
// @Description Get board columns by workspace
// @Tags board_columns
// @Accept json
// @Produce json
// @Param workspace_id path int true "Workspace ID"
// @Success 200 {array} models.TwBoardColumn
// @Router /dbms/v1/workspace/{workspace_id}/board_columns [get]
func (h *BoardColumnsHandler) getBoardColumnsByWorkspace(c *fiber.Ctx) error {
	// Parse the request
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	var boardColumns []models.TwBoardColumn
	// Get the board columns
	if result := h.DB.Where("workspace_id = ?", workspaceID).
		Where("deleted_at IS NULL").
		Order("position").
		Find(&boardColumns); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if boardColumns == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get board columns",
		})
	}
	// Return the response
	return c.JSON(boardColumns)
}

// getBoardColumnById godoc
// @Summary Get board column by ID
// @Description Get board column by ID
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id} [get]
func (h *BoardColumnsHandler) getBoardColumnById(c *fiber.Ctx) error {
	var boardColumn models.TwBoardColumn
	boardColumnId := c.Params("board_column_id")

	if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(boardColumn)
}

// deleteBoardColumn godoc
// @Summary Delete board column
// @Description Delete board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Success 204
// @Router /dbms/v1/board_columns/{id} [delete]
func (h *BoardColumnsHandler) deleteBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("board_column_id")
	var boardColumn models.TwBoardColumn

	// Retrieve the board column
	if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Update the deleted_at field using gorm.Expr("NOW()")
	if err := h.DB.Model(&boardColumn).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// updateBoardColumn godoc
// @Summary Update board column
// @Description Update board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param body body board_columns_dtos.BoardColumnsRequest true "Update board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id} [put]
func (h *BoardColumnsHandler) updateBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("board_column_id")
	var boardColumn models.TwBoardColumn
	if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	// Parse the request
	var updatedBoardColumn models.TwBoardColumn
	if err := c.BodyParser(&updatedBoardColumn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := h.DB.Model(&boardColumn).
		Updates(map[string]interface{}{
			"name":       updatedBoardColumn.Name,
			"updated_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(boardColumn)
}

// getBoardColumnField godoc
// @Summary Get board column field
// @Description Get board column field
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param field path string true "Field"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id}/{field} [get]
func (h *BoardColumnsHandler) getBoardColumnField(c *fiber.Ctx) error {
	field := c.Params("field")
	boardColumnId := c.Params("board_column_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	if field == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field != "name" && field != "position" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field == "name" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"name": boardColumn.Name,
		})
	}
	if field == "position" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(fiber.Map{
			"position": boardColumn.Position,
		})
	}
	return c.Status(fiber.StatusInternalServerError).SendString("Invalid field")
}

// updateBoardColumnField godoc
// @Summary Update board column field
// @Description Update board column field
// @Tags board_columns
// @Accept json
// @Produce json
// @Param id path int true "Board column ID"
// @Param field path string true "Field"
// @Param body body models.TwBoardColumn true "Update board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns/{id}/{field} [put]
func (h *BoardColumnsHandler) updateBoardColumnField(c *fiber.Ctx) error {
	field := c.Params("field")
	boardColumnId := c.Params("board_column_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	if field == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field != "name" && field != "position" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid field",
		})
	}
	if field == "name" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		var updateBoardColumnRequest models.TwBoardColumn
		if err := c.BodyParser(&updateBoardColumnRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if result := h.DB.Model(&boardColumn).Update("name", updateBoardColumnRequest.Name); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(updateBoardColumnRequest)

	}
	if field == "position" {
		var boardColumn models.TwBoardColumn
		if err := h.DB.Where("id = ?", boardColumnId).First(&boardColumn).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("BoardColumn not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		var updateBoardColumnRequest models.TwBoardColumn
		if err := c.BodyParser(&updateBoardColumnRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if result := h.DB.Model(&boardColumn).Update("position", updateBoardColumnRequest.Position); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
		return c.JSON(fiber.Map{
			"position": boardColumn.Position,
		})
	}
	return c.Status(fiber.StatusInternalServerError).SendString("Invalid field")
}

// createBoardColumn godoc
// @Summary Create board column
// @Description Create board column
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body board_columns_dtos.BoardColumnsRequest true "Create board column request"
// @Success 200 {object} models.TwBoardColumn
// @Router /dbms/v1/board_columns [post]
func (h *BoardColumnsHandler) createBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	var createBoardColumnRequest board_columns_dtos.BoardColumnsRequest
	if err := c.BodyParser(&createBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var boardColumn = models.TwBoardColumn{
		Name:        createBoardColumnRequest.Name,
		Position:    createBoardColumnRequest.Position,
		WorkspaceId: createBoardColumnRequest.WorkspaceId,
	}
	// Create the board column
	if result := h.DB.Create(&boardColumn); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(boardColumn)
}

func (h *BoardColumnsHandler) GetSchedulesByBoardColumn(c *fiber.Ctx) error {
	boardColumnId := c.Params("board_column_id")
	workspaceId := c.Params("workspace_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("board_column_id = ? and workspace_id =?", boardColumnId, workspaceId).Order("created_at").Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if schedules == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedules",
		})
	}
	return c.JSON(schedules)
}

type RequestBody struct {
	Position    int `json:"position"`
	WorkspaceId int `json:"workspace_id"`
}

// updatePositionAfterDeletion godoc
// @Summary Update position after deletion
// @Description Update position after deletion
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body RequestBody true "Update position after deletion request"
// @Success 200
// @Router /dbms/v1/board_columns/update_position_after_deletion [put]
func (h *BoardColumnsHandler) updatePositionAfterDeletion(c *fiber.Ctx) error {
	// Giải mã body request
	var requestBody RequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Lấy các tham số từ body
	position := requestBody.Position
	workspaceId := requestBody.WorkspaceId
	if position == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid position",
		})
	}
	if workspaceId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	// Update query to decrement the position of all columns with position greater than the specified position in the same workspace
	err := h.DB.Model(&models.TwBoardColumn{}).
		Where("position > ? AND workspace_id = ?", position, workspaceId).
		Where("deleted_at IS NULL").
		UpdateColumns(map[string]interface{}{
			"position":   gorm.Expr("position - 1"),
			"updated_at": gorm.Expr("NOW()"), // Cập nhật trường updated_at
		}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

type RageRequest struct {
	Position1   int `json:"position1"`
	Position2   int `json:"position2"`
	WorkspaceId int `json:"workspace_id"`
}

// getRage godoc
// @Summary Get rage
// @Description Get rage
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body RageRequest true "Get rage request"
// @Success 200 {array} models.TwBoardColumn
// @Router /dbms/v1/board_columns/rage/position [get]
func (h *BoardColumnsHandler) getRage(c *fiber.Ctx) error {
	// Giải mã body request
	var requestBody RageRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Lấy các tham số từ body
	position1 := requestBody.Position1
	position2 := requestBody.Position2
	workspaceId := requestBody.WorkspaceId

	// Kiểm tra tính hợp lệ của các tham số
	if position1 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid position1",
		})
	}
	if position2 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid position2",
		})
	}
	if workspaceId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}

	// Lấy các cột trong phạm vi vị trí và workspaceId
	var columns []models.TwBoardColumn
	if result := h.DB.Where("position >= ? AND position <= ? AND workspace_id = ?", position1, position2, workspaceId).
		Where("deleted_at IS NULL").
		Find(&columns); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Error.Error(),
		})
	}
	if columns == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get columns",
		})
	}

	// Trả về danh sách cột
	return c.JSON(columns)
}

// updatePosition godoc
// @Summary Update position
// @Description Update position
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body models.TwBoardColumn true "Update position request"
// @Success 200
// @Router /dbms/v1/board_columns/update_position/position [put]
func (h *BoardColumnsHandler) updatePosition(c *fiber.Ctx) error {

	var boardColumn models.TwBoardColumn
	if err := c.BodyParser(&boardColumn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	var oldBoardColumn models.TwBoardColumn
	if err := h.DB.Where("id = ?", boardColumn.ID).First(&oldBoardColumn).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to find the board column",
		})
	}
	err := h.DB.Model(&oldBoardColumn).UpdateColumns(map[string]interface{}{
		"position":   boardColumn.Position,
		"updated_at": gorm.Expr("NOW()"),
	}).Error
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Board column position updated successfully",
	})

}
