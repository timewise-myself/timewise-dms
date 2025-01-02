package recurrence_exception

import (
	"dbms/common"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRecurrenceExceptionHandler(router fiber.Router, db *gorm.DB) {
	recurrenceExceptionHandler := RecurrenceExceptionHandler{
		DB: db,
	}
	common.RegisterHandler(router, db, func(handler common.Handler) {
		handler.Router.Get("/", recurrenceExceptionHandler.GetRecurrenceExceptions)
		handler.Router.Get("/:schedule_id", recurrenceExceptionHandler.GetRecurrenceExceptionById)
		//handler.Router.Get("/user/:user_id", scheduleHandler.GetSchedulesByUserId)
		handler.Router.Post("/", recurrenceExceptionHandler.CreateRecurrenceException)
		handler.Router.Put("/:schedule_id", recurrenceExceptionHandler.UpdateRecurrenceException)
		handler.Router.Delete("/:schedule_id", recurrenceExceptionHandler.DeleteRecurrenceException)
	})
}
