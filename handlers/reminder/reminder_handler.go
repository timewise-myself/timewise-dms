package reminder

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type ReminderHandler struct {
	DB *gorm.DB
}

// createReminder godoc
// @Summary Create a reminder
// @Description Create a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder body models.TwReminder true "Reminder"
// @Success 200 {object} models.TwReminder
// @Router /dbms/v1/reminder [post]
func (h ReminderHandler) CreateReminder(ctx *fiber.Ctx) error {
	var reminder models.TwReminder
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if result := h.DB.Create(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

// getReminderById godoc
// @Summary Get reminder by ID
// @Description Get reminder by ID
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Success 200 {object} models.TwReminder
// @Router /dbms/v1/reminder/{reminder_id} [get]
func (h ReminderHandler) GetReminderById(ctx *fiber.Ctx) error {
	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

// getRemindersByScheduleId godoc
// @Summary Get reminders by schedule ID
// @Description Get reminders by schedule ID
// @Tags reminder
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwReminder
// @Router /dbms/v1/reminder/schedule/{schedule_id} [get]
func (h ReminderHandler) GetRemindersByScheduleId(ctx *fiber.Ctx) error {
	scheduleId := ctx.Params("schedule_id")
	var reminders []models.TwReminder
	if result := h.DB.
		Where("schedule_id = ?", scheduleId).
		Where("deleted_at IS NULL").
		Find(&reminders); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminders)
}

// updateReminder godoc
// @Summary Update a reminder
// @Description Update a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Param reminder body models.TwReminder true "Reminder"
// @Success 200 {object} models.TwReminder
// @Router /dbms/v1/reminder/{reminder_id} [put]
func (h ReminderHandler) UpdateReminder(ctx *fiber.Ctx) error {
	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if result := h.DB.Model(reminder).Update("updated_at", gorm.Expr("NOW()")); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	if result := h.DB.Omit("deleted_at", "created_at").Save(reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminder)
}

// deleteReminder godoc
// @Summary Delete a reminder
// @Description Delete a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Success 200 {object} models.TwReminder
// @Router /dbms/v1/reminder/{reminder_id} [delete]
func (h ReminderHandler) DeleteReminder(ctx *fiber.Ctx) error {

	id := ctx.Params("reminder_id")

	var reminder models.TwReminder

	// Ensure DB is initialized
	if h.DB == nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Database connection is not initialized")
	}

	// Query the database for the reminder
	if result := h.DB.Where("id = ?", id).Where("deleted_at IS NULL").First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}

	// Ensure reminder is valid
	if reminder.ID == 0 { // Check if reminder was found
		return ctx.Status(fiber.StatusNotFound).SendString("Reminder not found")
	}

	// Update the deleted_at field
	if result := h.DB.Model(&reminder).Update("deleted_at", gorm.Expr("NOW()")); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder deleted successfully",
	})
}

// getReminders godoc
// @Summary Get all reminders
// @Description Get all reminders
// @Tags reminder
// @Accept json
// @Produce json
// @Success 200 {array} models.TwReminder
// @Router /dbms/v1/reminder [get]
func (h ReminderHandler) GetReminders(ctx *fiber.Ctx) error {
	var reminders []models.TwReminder
	if result := h.DB.
		Where("deleted_at IS NULL").
		Preload("WorkspaceUser").
		Preload("WorkspaceUser.Workspace").
		Preload("WorkspaceUser.UserEmail").
		Preload("WorkspaceUser.UserEmail.User").
		Preload("Schedule").
		Find(&reminders); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return ctx.JSON(reminders)
}

// completeReminder godoc
// @Summary Complete a reminder
// @Description Complete a reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param reminder_id path string true "Reminder ID"
// @Success 200 {object} models.TwReminder
// @Router /dbms/v1/reminder/{reminder_id}/is_sent [put]
func (h ReminderHandler) CompleteReminder(ctx *fiber.Ctx) error {
	id := ctx.Params("reminder_id")
	var reminder models.TwReminder
	if result := h.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reminder); result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}

	// Chỉ cập nhật các thuộc tính cần thiết
	updateFields := map[string]interface{}{
		"updated_at": gorm.Expr("NOW()"),
		"is_sent":    true,
	}

	if result := h.DB.Model(&reminder).Updates(updateFields); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return ctx.JSON(reminder)
}
