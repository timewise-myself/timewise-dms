package schedule_log

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ScheduleLogHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterScheduleLogHandler(router fiber.Router, db *gorm.DB) {
	scheduleLogHandler := ScheduleLogHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", scheduleLogHandler.getScheduleLogs)
	router.Get("/:id", scheduleLogHandler.getScheduleLogById)
	router.Get("/schedule/:scheduleId", scheduleLogHandler.getScheduleLogsByScheduleID)
	router.Post("/", scheduleLogHandler.createScheduleLog)
	router.Put("/:id", scheduleLogHandler.updateScheduleLog)
	router.Delete("/:id", scheduleLogHandler.deleteScheduleLog)
}
