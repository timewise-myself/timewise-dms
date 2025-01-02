package schedule_participant

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

// getScheduleParticipants godoc
// @Summary Get all schedule participants
// @Description Get all schedule participants
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Success 200 {array} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant [get]
func (h *ScheduleParticipantHandler) getScheduleParticipants(c *fiber.Ctx) error {
	var scheduleParticipants []models.TwScheduleParticipant
	if result := h.DB.Find(&scheduleParticipants); result.Error != nil {
		// handle error
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(scheduleParticipants)
}

// @Summary Get schedule participant by ID
// @Description Get schedule participant by ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [get]
func (h *ScheduleParticipantHandler) getScheduleParticipantById(c *fiber.Ctx) error {
	var scheduleParticipant models.TwScheduleParticipant
	scheduleParticipantId := c.Params("id")

	if err := h.DB.Where("id = ?", scheduleParticipantId).First(&scheduleParticipant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("ScheduleParticipant not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(scheduleParticipant)
}

// updateScheduleParticipant godoc
// @Summary Update schedule participant
// @Description Update schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Param schedule_participant body models.TwScheduleParticipant true "Schedule participant object"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [put]
func (h *ScheduleParticipantHandler) UpdateScheduleParticipant(c *fiber.Ctx) error {
	var participantDTO schedule_participant_dtos.ScheduleParticipantRequest
	if err := c.BodyParser(&participantDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var participant models.TwScheduleParticipant
	participantID := c.Params("id")

	// Tìm participant theo ID
	if err := h.DB.Where("id = ?", participantID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Participant not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Kiểm tra và cập nhật từng trường
	if participantDTO.Status != nil {
		participant.Status = *participantDTO.Status
	}
	if participantDTO.AssignAt != nil {
		participant.AssignAt = participantDTO.AssignAt
	}
	if participantDTO.AssignBy != nil {
		participant.AssignBy = *participantDTO.AssignBy
	}
	if participantDTO.ResponseTime != nil {
		participant.ResponseTime = participantDTO.ResponseTime
	}
	if participantDTO.InvitationSentAt != nil {
		participant.InvitationSentAt = participantDTO.InvitationSentAt
	}
	if participantDTO.InvitationStatus != nil {
		participant.InvitationStatus = *participantDTO.InvitationStatus
	}

	// Cập nhật timestamp
	now := time.Now()
	participant.UpdatedAt = now

	// Lưu participant đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&participant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Trả về participant đã cập nhật
	return c.JSON(schedule_participant_dtos.ScheduleParticipantResponse{
		ID:               participant.ID,
		ScheduleId:       participant.ScheduleId,
		WorkspaceUserId:  participant.WorkspaceUserId,
		Status:           participant.Status,
		AssignAt:         participant.AssignAt,
		AssignBy:         participant.AssignBy,
		ResponseTime:     participant.ResponseTime,
		InvitationSentAt: participant.InvitationSentAt,
		InvitationStatus: participant.InvitationStatus,
	})
}

// deleteScheduleParticipant godoc
// @Summary Delete schedule participant
// @Description Delete schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param id path int true "Schedule Participant ID"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant/{id} [delete]
func (h *ScheduleParticipantHandler) deleteScheduleParticipant(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if result := h.DB.First(&scheduleParticipants, c.Params("id")); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	if result := h.DB.Delete(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(fiber.Map{
		"status": "deleted",
	})
}

// createScheduleParticipant godoc
// @Summary Create schedule participant
// @Description Create schedule participant
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param schedule_participant body models.TwScheduleParticipant true "Schedule participant object"
// @Success 200 {object} models.TwScheduleParticipant
// @Router /dbms/v1/schedule_participant [post]
func (h *ScheduleParticipantHandler) createScheduleParticipant(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if err := c.BodyParser(&scheduleParticipants); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	createSchedule := schedule_participant_dtos.ScheduleParticipantResponse{
		ID:               scheduleParticipants.ID,
		ScheduleId:       scheduleParticipants.ScheduleId,
		WorkspaceUserId:  scheduleParticipants.WorkspaceUserId,
		Status:           scheduleParticipants.Status,
		AssignAt:         scheduleParticipants.AssignAt,
		AssignBy:         scheduleParticipants.AssignBy,
		ResponseTime:     scheduleParticipants.ResponseTime,
		InvitationSentAt: scheduleParticipants.InvitationSentAt,
		InvitationStatus: scheduleParticipants.InvitationStatus,
	}
	return c.JSON(createSchedule)
}

func (h *ScheduleParticipantHandler) getScheduleParticipantByScheduleIdAndWorkspaceUserId(c *fiber.Ctx) error {
	var scheduleParticipant models.TwScheduleParticipant
	scheduleId := c.Params("scheduleId")
	workspaceUserId := c.Params("workspaceUserId")

	// Truy vấn database
	if err := h.DB.Where("workspace_user_id = ? AND schedule_id = ?", workspaceUserId, scheduleId).First(&scheduleParticipant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}

// isParticipantInSchedule godoc
// @Summary Check if workspace user is participant in the schedule
// @Description Check if any workspace user in the list is a participant in the schedule
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Param workspace_user_id query string true "List of workspace user IDs (comma-separated)"
// @Success 200 {boolean} true "User is participant"
// @Success 404 {boolean} false "User is not participant"
// @Router /dbms/v1/schedule_participant/{scheduleId}/participants [get]
func (h *ScheduleParticipantHandler) isParticipantInSchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("scheduleId")
	workspaceUserIds := c.Query("workspace_user_id") // Lấy danh sách workspaceUserIds từ query

	if workspaceUserIds == "" {
		return errors.New("workspaceUserIds parameter is required")
	}

	// Chuyển danh sách workspaceUserIds từ chuỗi thành mảng
	workspaceUserIdsList := strings.Split(workspaceUserIds, ",")

	// Kiểm tra xem có workspaceUserId nào thuộc scheduleId không
	var scheduleParticipants []models.TwScheduleParticipant
	if err := h.DB.Debug().Where("workspace_user_id IN (?) AND schedule_id = ? AND invitation_status = 'joined'", workspaceUserIdsList, scheduleId).Find(&scheduleParticipants).Error; err != nil {
		return err
	}

	// Nếu không tìm thấy participants nào, trả về false
	if len(scheduleParticipants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(errors.New("User is not participant"))
	}

	// Nếu tìm thấy ít nhất một participant, trả về true
	return c.Status(fiber.StatusOK).JSON(scheduleParticipants)
}

// getScheduleParticipantsByScheduleId godoc
// @Summary Get schedule participants by schedule ID
// @Description Get schedule participants by schedule ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Param workspaceId path string true "Workspace ID"
// @Success 200 {array} schedule_participant_dtos.ScheduleParticipantInfo
// @Router /dbms/v1/schedule_participant/workspace/{workspaceId}/schedule/{scheduleId} [get]
func (h *ScheduleParticipantHandler) getScheduleParticipantsByScheduleId(c *fiber.Ctx) error {
	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	scheduleId := c.Params("scheduleId")

	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schedule ID không hợp lệ",
		})
	}
	workspaceId := c.Params("workspaceId")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace ID không hợp lệ",
		})
	}

	// Perform the SQL query with multiple joins
	err := h.DB.Table("tw_schedule_participants AS sp").
		Select(`
            sp.id AS id,
            sp.schedule_id,
            sp.workspace_user_id,
			sp.status,
			sp.assign_at,
			sp.assign_by,
			sp.response_time,
			sp.invitation_sent_at,
			sp.invitation_status,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =sp.workspace_user_id").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("sp.schedule_id = ?", scheduleId).
		Where("wu.workspace_id = ?", workspaceId).
		Where("sp.deleted_at IS NULL").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("wu.is_active = true AND wu.is_verified = true AND wu.status = 'joined'").
		Scan(&scheduleParticipants).Error

	if err != nil {
		log.Println("Error querying schedule participants:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách participant",
		})
	}

	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy workspace",
		})
	}

	return c.JSON(scheduleParticipants)
}

// getScheduleParticipantsByScheduleId godoc
// @Summary Get schedule participants by schedule ID
// @Description Get schedule participants by schedule ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Param workspaceId path string true "Workspace ID"
// @Success 200 {array} schedule_participant_dtos.ScheduleParticipantInfo
// @Router /dbms/v1/schedule_participant/schedule/{scheduleId} [get]
func (h *ScheduleParticipantHandler) getScheduleParticipantsBySchedule(c *fiber.Ctx) error {
	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	scheduleId := c.Params("scheduleId")

	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schedule ID không hợp lệ",
		})
	}

	// Perform the SQL query with multiple joins
	err := h.DB.Table("tw_schedule_participants AS sp").
		Select(`
            sp.id AS id,
            sp.schedule_id,
            sp.workspace_user_id,
			sp.status,
			sp.assign_at,
			sp.assign_by,
			sp.response_time,
			sp.invitation_sent_at,
			sp.invitation_status,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =sp.workspace_user_id").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("sp.schedule_id = ?", scheduleId).
		Where("sp.deleted_at IS NULL AND sp.invitation_status != 'removed'").
		Where("wu.deleted_at IS NULL").
		Where("ue.deleted_at IS NULL").
		Where("u.deleted_at IS NULL").
		Where("wu.is_active = true AND wu.is_verified = true AND wu.status = 'joined'").
		Scan(&scheduleParticipants).Error

	if err != nil {
		log.Println("Error querying schedule participants:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách participant",
		})
	}

	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy workspace",
		})
	}

	return c.JSON(scheduleParticipants)
}

func (h *ScheduleParticipantHandler) inviteToSchedule(c *fiber.Ctx) error {
	var scheduleParticipants models.TwScheduleParticipant
	if err := c.BodyParser(&scheduleParticipants); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if result := h.DB.Create(&scheduleParticipants); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}
	return c.JSON(scheduleParticipants)
}

func (h *ScheduleParticipantHandler) RemoveParticipant(c *fiber.Ctx) error {

	var participant models.TwScheduleParticipant
	participantID := c.Params("id")

	// Tìm participant theo ID
	if err := h.DB.Where("id = ?", participantID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Participant not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	participant.Status = ""
	participant.InvitationStatus = "removed"

	// Cập nhật timestamp
	now := time.Now()
	participant.UpdatedAt = now

	// Lưu participant đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&participant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(schedule_participant_dtos.ScheduleParticipantResponse{
		ID:               participant.ID,
		ScheduleId:       participant.ScheduleId,
		WorkspaceUserId:  participant.WorkspaceUserId,
		Status:           participant.Status,
		AssignAt:         participant.AssignAt,
		AssignBy:         participant.AssignBy,
		ResponseTime:     participant.ResponseTime,
		InvitationSentAt: participant.InvitationSentAt,
		InvitationStatus: participant.InvitationStatus,
	})
}

func (h *ScheduleParticipantHandler) UnassignMember(c *fiber.Ctx) error {

	var participant models.TwScheduleParticipant
	participantID := c.Params("id")

	// Tìm participant theo ID
	if err := h.DB.Where("id = ?", participantID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Participant not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	participant.Status = "participant"

	// Cập nhật timestamp
	now := time.Now()
	participant.UpdatedAt = now

	// Lưu participant đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&participant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.JSON(schedule_participant_dtos.ScheduleParticipantResponse{
		ID:               participant.ID,
		ScheduleId:       participant.ScheduleId,
		WorkspaceUserId:  participant.WorkspaceUserId,
		Status:           participant.Status,
		AssignAt:         participant.AssignAt,
		AssignBy:         participant.AssignBy,
		ResponseTime:     participant.ResponseTime,
		InvitationSentAt: participant.InvitationSentAt,
		InvitationStatus: participant.InvitationStatus,
	})
}
