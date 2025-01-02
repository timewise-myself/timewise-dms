package recurrence_exception

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
)

type RecurrenceExceptionHandler struct {
	DB *gorm.DB
}

// GetRecurrenceExceptions godoc
// @Summary Get all recurrence exceptions
// @Description Get all recurrence exceptions
// @Tags recurrence_exception
// @Accept json
// @Produce json
// @Success 200 {array} core_dtos.TwRecurrenceExceptionResponseDTO
// @Router /dbms/v1/recurrence_exception [get]
func (h *RecurrenceExceptionHandler) GetRecurrenceExceptions(c *fiber.Ctx) error {
	var recurrenceExceptions []models.TwRecurrenceException
	if result := h.DB.Find(&recurrenceExceptions); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	var recurrenceExceptionDTOs []core_dtos.TwRecurrenceExceptionResponseDTO
	for _, recurrenceException := range recurrenceExceptions {
		recurrenceExceptionDTOs = append(recurrenceExceptionDTOs, core_dtos.TwRecurrenceExceptionResponseDTO{
			ID:            recurrenceException.ID,
			ScheduleId:    recurrenceException.ScheduleId,
			ExceptionDate: recurrenceException.ExceptionDate,
			NewStartTime:  recurrenceException.NewStartTime,
			NewEndTime:    recurrenceException.NewEndTime,
			IsCancelled:   recurrenceException.IsCancelled,
			ExtraData:     recurrenceException.ExtraData,
		})
	}
	return c.JSON(recurrenceExceptionDTOs)
}

// GetRecurrenceExceptionById godoc
// @Summary Get recurrence exception by ID
// @Description Get recurrence exception by ID
// @Tags recurrence_exception
// @Accept json
// @Produce json
// @Param id path int true "Recurrence Exception ID"
// @Success 200 {object} core_dtos.TwRecurrenceExceptionResponseDTO
// @Router /dbms/v1/recurrence_exception/{id} [get]
func (h *RecurrenceExceptionHandler) GetRecurrenceExceptionById(c *fiber.Ctx) error {
	id := c.Params("id")
	var recurrenceException models.TwRecurrenceException
	if result := h.DB.First(&recurrenceException, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	return c.JSON(core_dtos.TwRecurrenceExceptionResponseDTO{
		ID:            recurrenceException.ID,
		ScheduleId:    recurrenceException.ScheduleId,
		ExceptionDate: recurrenceException.ExceptionDate,
		NewStartTime:  recurrenceException.NewStartTime,
		NewEndTime:    recurrenceException.NewEndTime,
		IsCancelled:   recurrenceException.IsCancelled,
		ExtraData:     recurrenceException.ExtraData,
	})
}

// CreateRecurrenceException godoc
// @Summary Create a recurrence exception
// @Description Create a recurrence exception
// @Tags recurrence_exception
// @Accept json
// @Produce json
// @Param recurrence_exception body core_dtos.TwRecurrenceExceptionCreateDTO true "Recurrence Exception"
// @Success 200 {object} core_dtos.TwRecurrenceExceptionResponseDTO
// @Router /dbms/v1/recurrence_exception [post]
func (h *RecurrenceExceptionHandler) CreateRecurrenceException(c *fiber.Ctx) error {
	var recurrenceExceptionDTO core_dtos.TwRecurrenceExceptionCreateDTO
	if err := c.BodyParser(&recurrenceExceptionDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	recurrenceException := models.TwRecurrenceException{
		ScheduleId:    recurrenceExceptionDTO.ScheduleId,
		ExceptionDate: recurrenceExceptionDTO.ExceptionDate,
		NewStartTime:  recurrenceExceptionDTO.NewStartTime,
		NewEndTime:    recurrenceExceptionDTO.NewEndTime,
		IsCancelled:   recurrenceExceptionDTO.IsCancelled,
		ExtraData:     recurrenceExceptionDTO.ExtraData,
	}
	if result := h.DB.Create(&recurrenceException); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(core_dtos.TwRecurrenceExceptionResponseDTO{
		ID:            recurrenceException.ID,
		ScheduleId:    recurrenceException.ScheduleId,
		ExceptionDate: recurrenceException.ExceptionDate,
		NewStartTime:  recurrenceException.NewStartTime,
		NewEndTime:    recurrenceException.NewEndTime,
		IsCancelled:   recurrenceException.IsCancelled,
		ExtraData:     recurrenceException.ExtraData,
	})
}

// UpdateRecurrenceException godoc
// @Summary Update a recurrence exception
// @Description Update a recurrence exception
// @Tags recurrence_exception
// @Accept json
// @Produce json
// @Param id path int true "Recurrence Exception ID"
// @Param recurrence_exception body core_dtos.TwRecurrenceExceptionUpdateDTO true "Recurrence Exception"
// @Success 200 {object} core_dtos.TwRecurrenceExceptionResponseDTO
// @Router /dbms/v1/recurrence_exception/{id} [put]
func (h *RecurrenceExceptionHandler) UpdateRecurrenceException(c *fiber.Ctx) error {
	id := c.Params("id")
	var recurrenceException models.TwRecurrenceException
	if result := h.DB.First(&recurrenceException, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	var recurrenceExceptionDTO core_dtos.TwRecurrenceExceptionUpdateDTO
	if err := c.BodyParser(&recurrenceExceptionDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if recurrenceExceptionDTO.ExceptionDate != nil {
		recurrenceException.ExceptionDate = *recurrenceExceptionDTO.ExceptionDate
	}
	if recurrenceExceptionDTO.NewStartTime != nil {
		recurrenceException.NewStartTime = *recurrenceExceptionDTO.NewStartTime
	}
	if recurrenceExceptionDTO.NewEndTime != nil {
		recurrenceException.NewEndTime = *recurrenceExceptionDTO.NewEndTime
	}
	if recurrenceExceptionDTO.IsCancelled != nil {
		recurrenceException.IsCancelled = *recurrenceExceptionDTO.IsCancelled
	}
	if recurrenceExceptionDTO.ExtraData != nil {
		recurrenceException.ExtraData = *recurrenceExceptionDTO.ExtraData
	}
	if result := h.DB.Save(&recurrenceException); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(core_dtos.TwRecurrenceExceptionResponseDTO{
		ID:            recurrenceException.ID,
		ScheduleId:    recurrenceException.ScheduleId,
		ExceptionDate: recurrenceException.ExceptionDate,
		NewStartTime:  recurrenceException.NewStartTime,
		NewEndTime:    recurrenceException.NewEndTime,
		IsCancelled:   recurrenceException.IsCancelled,
		ExtraData:     recurrenceException.ExtraData,
	})
}

// DeleteRecurrenceException godoc
// @Summary Delete a recurrence exception
// @Description Delete a recurrence exception
// @Tags recurrence_exception
// @Accept json
// @Produce json
// @Param id path int true "Recurrence Exception ID"
// @Success 204
// @Router /dbms/v1/recurrence_exception/{id} [delete]
func (h *RecurrenceExceptionHandler) DeleteRecurrenceException(c *fiber.Ctx) error {
	id := c.Params("id")
	var recurrenceException models.TwRecurrenceException
	if result := h.DB.First(&recurrenceException, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString(result.Error.Error())
	}
	if result := h.DB.Delete(&recurrenceException, id); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
