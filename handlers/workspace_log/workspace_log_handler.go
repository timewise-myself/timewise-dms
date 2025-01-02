package workspace_log

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type WorkspaceLog struct {
	Router fiber.Router
	DB     *gorm.DB
}

// createWorkspaceLog godoc
// @Summary Create workspace log
// @Description Create workspace log
// @Tags workspace_log
// @Accept json
// @Produce json
// @Param workspace_log body models.TwWorkspaceLog true "Workspace log object"
// @Success 200 {object} models.TwWorkspaceLog
// @Router /dbms/v1/workspace_log [post]
func (h *WorkspaceLog) createWorkspaceLog(c *fiber.Ctx) error {
	var workspaceLog models.TwWorkspaceLog
	if err := c.BodyParser(&workspaceLog); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&workspaceLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(workspaceLog)
}

// @Summary Get all workspace logs
// @Description Get all workspace logs
// @Tags workspace_log
// @Accept json
// @Produce json
// @Success 200 {array} models.TwWorkspaceLog
// @Router /dbms/v1/workspace_log [get]
func (h *WorkspaceLog) getWorkspaceLog(c *fiber.Ctx) error {
	var workspaceLogs []models.TwWorkspaceLog
	if result := h.DB.Find(&workspaceLogs); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(workspaceLogs)
}

// @Summary Get workspace log by ID
// @Description Get workspace log by ID
// @Tags workspace_log
// @Accept json
// @Produce json
// @Param workspace_log_id path int true "Workspace log ID"
// @Success 200 {object} models.TwWorkspaceLog
// @Router /dbms/v1/workspace_log/{workspace_log_id} [get]
func (h *WorkspaceLog) getWorkspaceLogById(c *fiber.Ctx) error {
	var workspaceLog models.TwWorkspaceLog
	workspaceLogId := c.Params("workspace_log_id")

	if err := h.DB.Where("id = ?", workspaceLogId).First(&workspaceLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("WorkspaceLog not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(workspaceLog)
}

// removeWorkspaceLogById godoc
// @Summary Remove workspace log by ID
// @Description Remove workspace log by ID
// @Tags workspace_log
// @Accept json
// @Produce json
// @Param workspace_log_id path int true "Workspace log ID"
// @Success 204
// @Router /dbms/v1/workspace_log/{workspace_log_id} [delete]
func (h *WorkspaceLog) removeWorkspaceLogById(c *fiber.Ctx) error {
	workspaceLogId := c.Params("workspace_log_id")
	if result := h.DB.Delete(&models.TwWorkspaceLog{}, workspaceLogId); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)

}

func (h *WorkspaceLog) getWorkspaceLogsByWorkspaceId(ctx *fiber.Ctx) error {
	workspaceId := ctx.Params("workspace_id")
	var workspaceLogs []models.TwWorkspaceLog
	if result := h.DB.Where("workspace_id = ? and deleted_at IS NULL", workspaceId).Find(&workspaceLogs); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(workspaceLogs)

}
