package feature

import (
	_ "dbms/docs"
	"dbms/handlers/auth"
	"dbms/handlers/board_columns"
	comments "dbms/handlers/comments"
	"dbms/handlers/document"
	"dbms/handlers/notification"
	"dbms/handlers/notification_setting"
	"dbms/handlers/recurrence_exception"
	"dbms/handlers/reminder"
	"dbms/handlers/schedule"
	"dbms/handlers/schedule_log"
	"dbms/handlers/schedule_participant"
	"dbms/handlers/user"
	"dbms/handlers/user_email"
	"dbms/handlers/workspace"
	"dbms/handlers/workspace_log"
	"dbms/handlers/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

// @host localhost:8080
// @BasePath /dbms/v1
func RegisterHandlerV1(db *gorm.DB) *fiber.App {
	router := fiber.New()
	v1 := router.Group("/dbms/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	user.RegisterUserHandler(v1.Group("/user"), db)
	schedule_log.RegisterScheduleLogHandler(v1.Group("/schedule_log"), db)
	schedule_participant.RegisterScheduleParticipantHandler(v1.Group("/schedule_participant"), db)
	schedule.RegisterScheduleHandler(v1.Group("/schedule"), db)
	recurrence_exception.RegisterRecurrenceExceptionHandler(v1.Group("/recurrence_exception"), db)
	workspace_user.RegisterWorkspaceUserHandler(v1.Group("/workspace_user"), db)
	workspace_log.RegisterWorkspaceLogHandler(v1.Group("/workspace_log"), db)
	auth.RegisterAuthHandler(v1.Group("/auth"), db)
	user_email.RegisterUserEmailHandler(v1.Group("/user_email"), db)
	workspace.RegisterWorkspaceHandler(v1.Group("/workspace"), db)
	board_columns.RegisterBoardColumnsHandler(v1.Group("/board_columns"), db)
	document.RegisterDocumentHandler(v1.Group("/document"), db)
	comments.RegisterCommentsHandler(v1.Group("/comment"), db)
	notification.RegisterNotificationHandler(v1.Group("/notification"), db)
	reminder.RegisterReminderHandler(v1.Group("/reminder"), db)
	notification_setting.RegisterNotificationSettingHandler(v1.Group("/notification_setting"), db)
	return router
}
