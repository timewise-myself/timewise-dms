package schedule_log

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_log_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"log"
)

// getScheduleLogs godoc
// @Summary Get all schedule logs
// @Description Get all schedule logs
// @Tags schedule_log
// @Accept json
// @Produce json
// @Success 200 {array} models.TwScheduleLog
// @Router /dbms/v1/schedule_log [get]
func (h *ScheduleLogHandler) getScheduleLogs(c *fiber.Ctx) error {
	var scheduleLogs []models.TwScheduleLog
	if result := h.DB.Find(&scheduleLogs); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(scheduleLogs)
}

// @Summary Get schedule log by ID
// @Description Get schedule log by ID
// @Tags schedule_log
// @Accept json
// @Produce json
// @Param id path int true "Schedule Log ID"
// @Success 200 {object} models.TwScheduleLog
// @Router /dbms/v1/schedule_log/{id} [get]
func (h *ScheduleLogHandler) getScheduleLogById(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	scheduleLogId := c.Params("id")

	if err := h.DB.Where("id = ?", scheduleLogId).First(&scheduleLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("ScheduleLog not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(scheduleLog)
}

// updateScheduleLog godoc
// @Summary Update schedule log
// @Description Update schedule log
// @Tags schedule_log
// @Accept json
// @Produce json
// @Param id path int true "Schedule Log ID"
// @Param schedule_log body models.TwScheduleLog true "Schedule log object"
// @Success 200 {object} models.TwScheduleLog
// @Router /dbms/v1/schedule_log/{id} [put]
func (h *ScheduleLogHandler) updateScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog

	if result := h.DB.First(&scheduleLog, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := h.DB.Save(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleLog)
}

// deleteScheduleLog godoc
// @Summary Delete schedule log
// @Description Delete schedule log
// @Tags schedule_log
// @Accept json
// @Produce json
// @Param id path int true "Schedule Log ID"
// @Success 200 {object} models.TwScheduleLog
// @Router /dbms/v1/schedule_log/{id} [delete]
func (h *ScheduleLogHandler) deleteScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	if result := h.DB.First(&scheduleLog, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	if result := h.DB.Delete(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

// createScheduleLog godoc
// @Summary Create a new schedule log
// @Description Create a new schedule log
// @Tags schedule_log
// @Accept json
// @Produce json
// @Param schedule_log body models.TwScheduleLog true "Schedule log object"
// @Success 200 {object} models.TwScheduleLog
// @Router /dbms/v1/schedule_log [post]
func (h *ScheduleLogHandler) createScheduleLog(c *fiber.Ctx) error {
	var scheduleLog models.TwScheduleLog
	if err := c.BodyParser(&scheduleLog); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&scheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleLog)
}

func (h *ScheduleLogHandler) getScheduleLogsByScheduleID(c *fiber.Ctx) error {
	var scheduleLogs []schedule_log_dtos.TwScheduleLogResponse
	scheduleId := c.Params("scheduleId")

	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schedule ID không hợp lệ",
		})
	}

	// Perform the SQL query with multiple joins
	err := h.DB.Table("tw_schedule_logs AS sl").
		Select(`
            sl.id AS id,
            sl.created_at,
            sl.updated_at,
			sl.schedule_id,
			sl.workspace_user_id,
			sl.action,
			sl.field_changed,
			sl.old_value,
			sl.new_value,
			sl.description,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =sl.workspace_user_id").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("sl.schedule_id = ?", scheduleId).
		Where("sl.deleted_at IS NULL").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("wu.is_active = true AND wu.is_verified = true AND wu.status = 'joined'").
		Scan(&scheduleLogs).Error

	if err != nil {
		log.Println("Error querying schedule participants:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách participant",
		})
	}

	return c.JSON(scheduleLogs)
}
