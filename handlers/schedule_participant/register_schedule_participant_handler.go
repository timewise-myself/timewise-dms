package schedule_participant

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ScheduleParticipantHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterScheduleParticipantHandler(router fiber.Router, db *gorm.DB) {
	scheduleParticipantHandeler := ScheduleParticipantHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/", scheduleParticipantHandeler.getScheduleParticipants)
	router.Get("/:id", scheduleParticipantHandeler.getScheduleParticipantById)
	router.Get("/:scheduleId/participants", scheduleParticipantHandeler.isParticipantInSchedule)
	router.Get("/workspace_user/:workspaceUserId/schedule/:scheduleId", scheduleParticipantHandeler.getScheduleParticipantByScheduleIdAndWorkspaceUserId)
	router.Post("/", scheduleParticipantHandeler.createScheduleParticipant)
	router.Put("/:id", scheduleParticipantHandeler.UpdateScheduleParticipant)
	router.Delete("/:id", scheduleParticipantHandeler.deleteScheduleParticipant)
	router.Get("/workspace/:workspaceId/schedule/:scheduleId", scheduleParticipantHandeler.getScheduleParticipantsByScheduleId)
	router.Get("/schedule/:scheduleId", scheduleParticipantHandeler.getScheduleParticipantsBySchedule)
	router.Post("/invite", scheduleParticipantHandeler.inviteToSchedule)
	router.Put("/remove/:id", scheduleParticipantHandeler.RemoveParticipant)
	router.Put("/unassign/:id", scheduleParticipantHandeler.UnassignMember)
}
