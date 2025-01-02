package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type ScheduleHandler struct {
	DB *gorm.DB
}

func parseTime(timeStr string) (time.Time, error) {
	// Sử dụng thư viện dateparse để phân tích chuỗi thời gian
	layout := "2006-01-02 15:04:05.000"
	parsedTime, err := time.ParseInLocation(layout, timeStr, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime.UTC(), nil
}

// FilterSchedules godoc
// @Summary Filter schedule
// @Description Filter schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param workspace_id query int false "Workspace ID"
// @Param board_column_id query int false "Board Column ID"
// @Param title query string false "Title of the schedule (searches with LIKE)"
// @Param start_time query string false "Start time of the schedule (ISO8601 format, filter by schedules starting after this date)"
// @Param end_time query string false "End time of the schedule (ISO8601 format, filter by schedules ending before this date)"
// @Param location query string false "Location of the schedule (searches with LIKE)"
// @Param created_by query int false "User ID of the creator"
// @Param status query string false "Status of the schedule"
// @Param is_deleted query bool false "Filter by deleted schedules"
// @Param assigned_to query int false "User ID assigned to the schedule"
// @Success 200 {array} core_dtos.TwScheduleResponse "Filtered list of schedules"
// @Failure 400 {object} fiber.Error "Invalid query parameters"
// @Failure 500 {object} fiber.Error "Internal Server Error"
// @Router /dbms/v1/schedule/schedules/filter [get]
func (h *ScheduleHandler) FilterSchedules(c *fiber.Ctx) error {
	var schedules []models.TwSchedule

	query := h.DB.Table("tw_schedules").
		Joins("JOIN tw_workspaces ON tw_schedules.workspace_id = tw_workspaces.id AND tw_workspaces.deleted_at IS NULL").
		Joins("JOIN tw_board_columns ON tw_schedules.board_column_id = tw_board_columns.id AND tw_board_columns.deleted_at IS NULL")
	workspaceID := c.Query("workspace_id")
	boardColumnID := c.Query("board_column_id")
	title := c.Query("title")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	location := c.Query("location")
	createdBy := c.Query("created_by")
	status := c.Query("status")
	isDeleted := c.Query("is_deleted")
	assignedTo := c.Query("assigned_to")

	if workspaceID != "" {
		workspaceIDSubStrings := strings.Split(workspaceID, ",")
		query = query.Where("tw_schedules.workspace_id IN (?)", workspaceIDSubStrings)
	}

	if boardColumnID != "" {
		query = query.Where("tw_schedules.board_column_id = ?", boardColumnID)
	}

	if title != "" {
		query = query.Where("tw_schedules.title LIKE ?", "%"+title+"%")
	}

	if startTime != "" {
		parsedStartTime, err := parseTime(startTime)
		if err != nil {
			return err
		}
		query = query.Where("tw_schedules.start_time >= ?", parsedStartTime)
	}

	if endTime != "" {
		parsedEndTime, err := parseTime(endTime)
		if err != nil {
			return err
		}
		query = query.Where("tw_schedules.end_time <= ?", parsedEndTime)
	}

	if location != "" {
		query = query.Where("tw_schedules.location LIKE ?", "%"+location+"%")
	}

	if createdBy != "" {
		query = query.Where("tw_schedules.created_by = ?", createdBy)
	}

	if status != "" {
		query = query.Where("tw_schedules.status = ?", status)
	}

	if isDeleted != "" {
		if isDeleted == "true" {
			query = query.Where("tw_schedules.is_deleted = ?", 1)
		} else if isDeleted == "false" {
			query = query.Where("tw_schedules.is_deleted = ?", 0)
		} else {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid value for is_deleted. Must be 'true' or 'false'")
		}
	}

	if assignedTo != "" {
		query = query.Where("tw_schedules.assigned_to @> ?", "{"+assignedTo+"}")
	}

	if result := query.Debug().Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	var scheduleDTOs []core_dtos.TwScheduleResponse
	for _, schedule := range schedules {
		scheduleDTO := core_dtos.TwScheduleResponse{
			ID:                int(schedule.ID),
			WorkspaceID:       schedule.WorkspaceId,
			BoardColumnID:     schedule.BoardColumnId,
			Title:             schedule.Title,
			Description:       schedule.Description,
			Location:          schedule.Location,
			CreatedBy:         schedule.CreatedBy,
			Status:            schedule.Status,
			AllDay:            schedule.AllDay,
			Visibility:        schedule.Visibility,
			ExtraData:         schedule.ExtraData,
			IsDeleted:         schedule.IsDeleted,
			RecurrencePattern: schedule.RecurrencePattern,
		}

		// Check if StartTime is not nil before dereferencing
		if schedule.StartTime != nil {
			scheduleDTO.StartTime = *schedule.StartTime
		}

		// Check if EndTime is not nil before dereferencing
		if schedule.EndTime != nil {
			scheduleDTO.EndTime = *schedule.EndTime
		}

		// Check if CreatedAt is not nil before dereferencing
		if schedule.CreatedAt != nil {
			scheduleDTO.CreatedAt = *schedule.CreatedAt
		}

		// Check if UpdatedAt is not nil before dereferencing
		if schedule.UpdatedAt != nil {
			scheduleDTO.UpdatedAt = *schedule.UpdatedAt
		}

		scheduleDTOs = append(scheduleDTOs, scheduleDTO)
	}

	return c.JSON(scheduleDTOs)
}

// GetSchedules godoc
// @Summary Get all schedules
// @Description Get all schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Success 200 {array} core_dtos.TwScheduleResponse
// @Router /dbms/v1/schedule [get]
func (h *ScheduleHandler) GetSchedules(c *fiber.Ctx) error {
	var schedules []models.TwSchedule
	if result := h.DB.Find(&schedules); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	var scheduleDTOs []core_dtos.TwScheduleResponse
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, core_dtos.TwScheduleResponse{
			ID:                int(schedule.ID),
			WorkspaceID:       schedule.WorkspaceId,
			BoardColumnID:     schedule.BoardColumnId,
			Title:             schedule.Title,
			Description:       schedule.Description,
			StartTime:         *schedule.StartTime,
			EndTime:           *schedule.EndTime,
			Location:          schedule.Location,
			CreatedBy:         schedule.CreatedBy,
			CreatedAt:         *schedule.CreatedAt,
			UpdatedAt:         *schedule.UpdatedAt,
			Status:            schedule.Status,
			AllDay:            schedule.AllDay,
			Visibility:        schedule.Visibility,
			ExtraData:         schedule.ExtraData,
			IsDeleted:         schedule.IsDeleted,
			RecurrencePattern: schedule.RecurrencePattern,
			//AssignedTo:        []int{schedule.AssignedTo},
		})
	}

	return c.JSON(scheduleDTOs)
}

// GetScheduleById godoc
// @Summary Get schedule by ID
// @Description Get schedule by ID
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} core_dtos.TwScheduleResponse
// @Router /dbms/v1/schedule/{schedule_id} [get]
func (h *ScheduleHandler) GetScheduleById(c *fiber.Ctx) error {
	var schedule models.TwSchedule
	scheduleId := c.Params("schedule_id")

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var startTime, endTime, createdAt, updatedAt time.Time

	if schedule.StartTime != nil {
		startTime = *schedule.StartTime
	}
	if schedule.EndTime != nil {
		endTime = *schedule.EndTime
	}
	if schedule.CreatedAt != nil {
		createdAt = *schedule.CreatedAt
	}
	if schedule.UpdatedAt != nil {
		updatedAt = *schedule.UpdatedAt
	}

	scheduleDTO := core_dtos.TwScheduleResponse{
		ID:                int(schedule.ID),
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         startTime,
		EndTime:           endTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		VideoTranscript:   &schedule.VideoTranscript,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		Position:          schedule.Position,
		Priority:          schedule.Priority,
		//AssignedTo:        []int{schedule.AssignedTo},
	}

	return c.JSON(scheduleDTO)
}

// CreateSchedule godoc
// @Summary Create a new schedule
// @Description Create a new schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule body core_dtos.TwCreateScheduleRequest true "Schedule"
// @Success 201 {object} core_dtos.TwCreateShecduleResponse
// @Router /dbms/v1/schedule [post]
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {

	var scheduleDTO core_dtos.TwCreateScheduleRequest
	if err := c.BodyParser(&scheduleDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var existingSchedules []models.TwSchedule
	if err := h.DB.Where("board_column_id = ? and is_deleted = false", *scheduleDTO.BoardColumnID).Find(&existingSchedules).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now()
	endTime := now.Add(1 * time.Hour)
	schedule := models.TwSchedule{
		WorkspaceId:   *scheduleDTO.WorkspaceID,
		BoardColumnId: *scheduleDTO.BoardColumnID,
		Title:         *scheduleDTO.Title,
		Description:   *scheduleDTO.Description,
		StartTime:     &now,
		EndTime:       &endTime,
		CreatedBy:     *scheduleDTO.WorkspaceUserID,
		CreatedAt:     &now,
		UpdatedAt:     &now,
		Position:      len(existingSchedules) + 1,
		Status:        "not yet",
		Visibility:    "public",
	}

	if scheduleDTO.StartTime != nil {
		isoTime := convertToISOFormat(*scheduleDTO.StartTime)
		parsedTime := convertDateFormat(&isoTime)
		if parsedTime != nil {
			schedule.StartTime = parsedTime
		}

	}

	if scheduleDTO.EndTime != nil {
		isoTime := convertToISOFormat(*scheduleDTO.EndTime)
		parsedTime := convertDateFormat(&isoTime)
		if parsedTime != nil {
			schedule.EndTime = parsedTime
		}
	}

	if result := h.DB.Create(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	newScheduleLog := models.TwScheduleLog{
		ScheduleId:      schedule.ID,
		WorkspaceUserId: *scheduleDTO.WorkspaceUserID,
		Action:          "create schedule",
	}

	if result := h.DB.Create(&newScheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	now = time.Now()
	newScheduleParticipant := models.TwScheduleParticipant{
		CreatedAt:        now,
		UpdatedAt:        now,
		ScheduleId:       schedule.ID,
		WorkspaceUserId:  *scheduleDTO.WorkspaceUserID,
		AssignAt:         &now,
		AssignBy:         *scheduleDTO.WorkspaceUserID,
		Status:           "creator",
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	if result := h.DB.Create(&newScheduleParticipant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(core_dtos.TwCreateShecduleResponse{
		ID:            schedule.ID,
		WorkspaceID:   schedule.WorkspaceId,
		BoardColumnID: schedule.BoardColumnId,
		Title:         schedule.Title,
		Description:   schedule.Description,
		Position:      schedule.Position,
		StartTime:     *schedule.StartTime,
		EndTime:       *schedule.EndTime,
	})
}

func convertToISOFormat(input string) string {
	// Thay dấu cách bằng 'T' và bỏ phần '.000' ở cuối, sau đó thêm 'Z'
	isoFormat := strings.Replace(input, " ", "T", 1)
	isoFormat = strings.TrimSuffix(isoFormat, ".000") + "Z"
	return isoFormat
}

func convertDateFormat(dateStr *string) *time.Time {
	if dateStr == nil {
		return nil
	}

	// Định dạng ngày giờ đầu vào, bao gồm múi giờ (UTC nếu không có múi giờ cụ thể)
	const inputFormat = "2006-01-02T15:04:05Z07:00"

	// Phân tích chuỗi ngày giờ theo múi giờ UTC
	parsedTime, err := time.Parse(inputFormat, *dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil
	}

	// Trả về giá trị UTC để lưu trữ
	utcTime := parsedTime.UTC()
	return &utcTime
}

// UpdateSchedule godoc
// @Summary Update an existing schedule
// @Description Update an existing schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Param schedule body core_dtos.TwUpdateScheduleRequest true "Schedule"
// @Success 200 {object} core_dtos.TwUpdateScheduleResponse
// @Router /dbms/v1/schedule/{schedule_id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	var scheduleDTO core_dtos.TwUpdateScheduleRequest
	if err := c.BodyParser(&scheduleDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var schedule models.TwSchedule

	scheduleId := c.Params("schedule_id")
	workspaceUserIdStr := c.Params("workspace_user_id")
	workspaceUserId, err := strconv.Atoi(workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid workspace_user_id")
	}

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Tạo danh sách các log khi trường được cập nhật
	var logs []models.TwScheduleLog

	// Hàm phụ: Kiểm tra và ghi log nếu có thay đổi
	checkAndLog := func(field, oldValue, newValue string) {
		if oldValue != newValue {
			logs = append(logs, models.TwScheduleLog{
				ScheduleId:      schedule.ID,
				WorkspaceUserId: workspaceUserId,
				Action:          "update schedule",
				FieldChanged:    field,
				OldValue:        oldValue,
				NewValue:        newValue,
			})
		}
	}

	if scheduleDTO.Title != nil {
		checkAndLog("title", schedule.Title, *scheduleDTO.Title)
		schedule.Title = *scheduleDTO.Title
	}
	if scheduleDTO.Description != nil {
		checkAndLog("description", schedule.Description, *scheduleDTO.Description)
		schedule.Description = *scheduleDTO.Description
	}
	if scheduleDTO.StartTime != nil {
		oldStartTime := ""
		if schedule.StartTime != nil {
			oldStartTime = schedule.StartTime.String()
		}
		checkAndLog("start_time", oldStartTime, *scheduleDTO.StartTime)
		parsedTime := convertDateFormat(scheduleDTO.StartTime)
		if parsedTime != nil {
			schedule.StartTime = parsedTime
		}

	}

	if scheduleDTO.EndTime != nil {
		oldEndTime := ""
		if schedule.EndTime != nil {
			oldEndTime = schedule.EndTime.String()
		}
		checkAndLog("end_time", oldEndTime, *scheduleDTO.EndTime)
		parsedTime := convertDateFormat(scheduleDTO.EndTime)
		if parsedTime != nil {
			schedule.EndTime = parsedTime
		}
	}

	if scheduleDTO.Location != nil {
		checkAndLog("location", schedule.Location, *scheduleDTO.Location)
		schedule.Location = *scheduleDTO.Location
	}
	if scheduleDTO.Status != nil {
		checkAndLog("status", schedule.Status, *scheduleDTO.Status)
		schedule.Status = *scheduleDTO.Status
	}
	if scheduleDTO.AllDay != nil {
		checkAndLog("all_day", strconv.FormatBool(schedule.AllDay), strconv.FormatBool(*scheduleDTO.AllDay))
		schedule.AllDay = *scheduleDTO.AllDay
	}
	if scheduleDTO.Visibility != nil {
		checkAndLog("visibility", schedule.Visibility, *scheduleDTO.Visibility)
		schedule.Visibility = *scheduleDTO.Visibility
	}
	if scheduleDTO.ExtraData != nil {
		checkAndLog("extra_data", schedule.ExtraData, *scheduleDTO.ExtraData)
		schedule.ExtraData = *scheduleDTO.ExtraData
	}
	if scheduleDTO.RecurrencePattern != nil {
		checkAndLog("recurrence_pattern", schedule.RecurrencePattern, *scheduleDTO.RecurrencePattern)
		schedule.RecurrencePattern = *scheduleDTO.RecurrencePattern
	}
	if scheduleDTO.Priority != nil {
		checkAndLog("priority", schedule.Priority, *scheduleDTO.Priority)
		schedule.Priority = *scheduleDTO.Priority
	}
	if scheduleDTO.VideoTranscript != nil {
		checkAndLog("video_transcript", schedule.VideoTranscript, *scheduleDTO.VideoTranscript)
		schedule.VideoTranscript = *scheduleDTO.VideoTranscript
	}
	schedule.CreatedBy = workspaceUserId

	// Update timestamp
	now := time.Now()
	schedule.UpdatedAt = &now

	// Lưu schedule đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Thêm các log vào cơ sở dữ liệu
	if len(logs) > 0 {
		if result := h.DB.Create(&logs); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
	}

	// Trả về kết quả cập nhật thành công
	return c.JSON(core_dtos.TwUpdateScheduleResponse{
		ID:                schedule.ID,
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         *schedule.StartTime,
		EndTime:           *schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         *schedule.CreatedAt,
		UpdatedAt:         *schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		Position:          schedule.Position,
		Priority:          schedule.Priority,
		VideoTranscript:   schedule.VideoTranscript,
	})
}

func (h *ScheduleHandler) UpdateSchedulePosition(c *fiber.Ctx) error {
	var scheduleDTO core_dtos.TwUpdateSchedulePosition
	if err := c.BodyParser(&scheduleDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var schedule models.TwSchedule

	scheduleId := c.Params("schedule_id")
	workspaceUserIdStr := c.Params("workspace_user_id")
	workspaceUserId, err := strconv.Atoi(workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid workspace_user_id")
	}

	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var logs []models.TwScheduleLog

	checkAndLog := func(field, oldValue, newValue string) {
		if oldValue != newValue {
			logs = append(logs, models.TwScheduleLog{
				ScheduleId:      schedule.ID,
				WorkspaceUserId: workspaceUserId,
				Action:          "update schedule",
				FieldChanged:    field,
				OldValue:        oldValue,
				NewValue:        newValue,
			})
		}
	}

	if scheduleDTO.BoardColumnID != nil {
		if *scheduleDTO.BoardColumnID == schedule.BoardColumnId {
			if *scheduleDTO.Position < schedule.Position {
				var schedulesToUpdate []models.TwSchedule
				if err := h.DB.Where("board_column_id = ? AND position < ? AND position >= ? AND is_deleted != 1", schedule.BoardColumnId, schedule.Position, scheduleDTO.Position).
					Order("position ASC").Find(&schedulesToUpdate).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}

				for i := range schedulesToUpdate {
					schedulesToUpdate[i].Position += 1
					if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesToUpdate[i]).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
					}
				}
			} else if *scheduleDTO.Position > schedule.Position {
				var schedulesToUpdate []models.TwSchedule
				if err := h.DB.Where("board_column_id = ? AND position > ? AND position <= ? AND is_deleted != 1", schedule.BoardColumnId, schedule.Position, scheduleDTO.Position).
					Order("position ASC").Find(&schedulesToUpdate).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}

				for i := range schedulesToUpdate {
					schedulesToUpdate[i].Position -= 1
					if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesToUpdate[i]).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
					}
				}
			}
			checkAndLog("position", strconv.Itoa(schedule.Position), strconv.Itoa(*scheduleDTO.Position))
			schedule.Position = *scheduleDTO.Position
		} else if *scheduleDTO.BoardColumnID != schedule.BoardColumnId {
			var schedulesToUpdate []models.TwSchedule
			if err := h.DB.Where("board_column_id = ? AND position > ? AND is_deleted != 1", schedule.BoardColumnId, schedule.Position).
				Order("position ASC").Find(&schedulesToUpdate).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			for i := range schedulesToUpdate {
				schedulesToUpdate[i].Position -= 1
				if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesToUpdate[i]).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}
			}
			var schedulesInBoardColumnToUpdate []models.TwSchedule
			if err := h.DB.Where("board_column_id = ? AND position >= ? AND is_deleted != 1", scheduleDTO.BoardColumnID, scheduleDTO.Position).
				Order("position ASC").Find(&schedulesInBoardColumnToUpdate).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			if len(schedulesInBoardColumnToUpdate) == 0 {
				var maxPosition int
				if err := h.DB.Model(&models.TwSchedule{}).
					Where("board_column_id = ? AND is_deleted != 1", scheduleDTO.BoardColumnID).
					Select("COALESCE(MAX(position), 0)").Scan(&maxPosition).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}
				checkAndLog("position", strconv.Itoa(schedule.Position), strconv.Itoa(*scheduleDTO.Position))
				schedule.Position = maxPosition + 1
			} else {
				for i := range schedulesInBoardColumnToUpdate {
					schedulesInBoardColumnToUpdate[i].Position += 1
					if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesInBoardColumnToUpdate[i]).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
					}
				}
				checkAndLog("position", strconv.Itoa(schedule.Position), strconv.Itoa(*scheduleDTO.Position))
				schedule.Position = *scheduleDTO.Position
			}
		}
		checkAndLog("board_column_id", strconv.Itoa(schedule.BoardColumnId), strconv.Itoa(*scheduleDTO.BoardColumnID))
		schedule.BoardColumnId = *scheduleDTO.BoardColumnID
	}

	// Update timestamp
	now := time.Now()
	schedule.UpdatedAt = &now

	// Lưu schedule đã cập nhật
	if result := h.DB.Omit("deleted_at").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Thêm các log vào cơ sở dữ liệu
	if len(logs) > 0 {
		if result := h.DB.Create(&logs); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
		}
	}

	// Trả về kết quả cập nhật thành công
	return c.JSON(core_dtos.TwUpdateScheduleResponse{
		ID:                schedule.ID,
		WorkspaceID:       schedule.WorkspaceId,
		BoardColumnID:     schedule.BoardColumnId,
		Title:             schedule.Title,
		Description:       schedule.Description,
		StartTime:         *schedule.StartTime,
		EndTime:           *schedule.EndTime,
		Location:          schedule.Location,
		CreatedBy:         schedule.CreatedBy,
		CreatedAt:         *schedule.CreatedAt,
		UpdatedAt:         *schedule.UpdatedAt,
		Status:            schedule.Status,
		AllDay:            schedule.AllDay,
		Visibility:        schedule.Visibility,
		ExtraData:         schedule.ExtraData,
		IsDeleted:         schedule.IsDeleted,
		RecurrencePattern: schedule.RecurrencePattern,
		Position:          schedule.Position,
		Priority:          schedule.Priority,
	})
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Description Delete a schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 204 "No Content"
// @Router /dbms/v1/schedule/{schedule_id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")
	workspaceUserIdStr := c.Params("workspace_user_id")
	workspaceUserId, err := strconv.Atoi(workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid workspace_user_id")
	}

	var schedule models.TwSchedule
	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now()

	schedule.IsDeleted = true
	schedule.UpdatedAt = &now
	schedule.DeletedAt = &now

	if result := h.DB.Omit("start_time,end_time").Save(&schedule); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	var schedulesToUpdate []models.TwSchedule
	if err := h.DB.Where("board_column_id = ? AND position > ? AND is_deleted != 1", schedule.BoardColumnId, schedule.Position).
		Order("position ASC").Find(&schedulesToUpdate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for i := range schedulesToUpdate {
		schedulesToUpdate[i].Position -= 1
		if err := h.DB.Omit("deleted_at,start_time,end_time").Save(&schedulesToUpdate[i]).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	newScheduleLog := models.TwScheduleLog{
		ScheduleId:      schedule.ID,
		WorkspaceUserId: workspaceUserId,
		Action:          "delete schedule",
	}

	if result := h.DB.Create(&newScheduleLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ScheduleHandler) GetSchedulesByBoardColumn(c *fiber.Ctx) error {
	boardColumnID := c.Params("board_column_id")
	if boardColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("board_column_id = ?", boardColumnID).Find(&schedules); result.Error != nil {
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

func (h *ScheduleHandler) GetSchedulesByWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("workspace_id = ?", workspaceID).Find(&schedules); result.Error != nil {
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

func (h *ScheduleHandler) getSchedulesByBoardColumn(c *fiber.Ctx) error {
	boardColumnID := c.Params("board_column_id")
	workspaceID := c.Params("workspace_id")
	if boardColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	var schedules []models.TwSchedule
	if result := h.DB.Where("board_column_id = ? and workspace_id = ? and is_deleted = false", boardColumnID, workspaceID).
		Order("position").
		Find(&schedules); result.Error != nil {
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

// UpdateTranscriptBySchedule godoc
// @Summary Update transcript by schedule
// @Description Update transcript by schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param x_api_key header string true "API Key"
// @Param schedule_id path string true "Schedule ID"
// @Param video_transcript formData string true "Video transcript"
// @Success 200 "Updated successfully"
// @Router /dbms/v1/schedule/{schedule_id}/transcript [put]
func (h *ScheduleHandler) UpdateTranscriptBySchedule(ctx *fiber.Ctx) error {
	// get an api_key from params
	//apiKey := ctx.Get("x_api_key")
	//if apiKey != "667qwsrUlyVa" {
	//	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//		"error": "Unauthorized",
	//	})
	//}

	// Parse schedule_id from params
	scheduleId := ctx.Params("schedule_id")
	if scheduleId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Schedule ID is required")
	}

	// Get video_transcript from body data
	bodyDataStr := ctx.BodyRaw()
	var bodyData map[string]interface{}
	_ = json.Unmarshal(bodyDataStr, &bodyData)
	videoTranscript, ok := bodyData["video_transcript"].(map[string]interface{})
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Video transcript is required")
	}
	if videoTranscript == nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Video transcript is required")
	}

	// Fetch the schedule from the database
	var schedule models.TwSchedule
	if err := h.DB.Where("id = ?", scheduleId).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString("Schedule not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Convert the JSON object to a string
	videoTranscriptStr, _ := json.Marshal(videoTranscript)

	schedule.VideoTranscript = string(videoTranscriptStr)

	now := time.Now()
	schedule.UpdatedAt = &now

	// Save the updated schedule back to the database
	if result := h.DB.Save(&schedule); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Return the updated schedule in the response
	return ctx.JSON("Updated successfully")
}

// GetSchedulesByBoardColumnFilter godoc
// @Summary Get schedules by board column with filters
// @Description Get schedules by board column with filters
// @Tags schedule
// @Accept json
// @Produce json
// @Param workspace_id path int true "Workspace ID"
// @Param board_column_id path int true "Board Column ID"
// @Param search query string false "Search by schedule title"
// @Param member query string false "Filter by member emails"
// @Param due query string false "Filter by due date (day, week, month)"
// @Param dueComplete query string false "Filter by due complete"
// @Param overdue query string false "Filter by overdue"
// @Param notDue query string false "Filter by not due"
// @Success 200 {array} models.TwSchedule
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /dbms/v1/schedule/workspace/{workspace_id}/board_column/{board_column_id}/filter [get]
func (h *ScheduleHandler) getSchedulesByBoardColumnFilter(c *fiber.Ctx) error {
	boardColumnID := c.Params("board_column_id")
	workspaceID := c.Params("workspace_id")
	if boardColumnID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid board column ID",
		})
	}
	search := c.Query("search", "")
	membersParam := c.Query("member", "")
	dueParam := c.Query("due")
	dueCompleteParam := c.Query("dueComplete")
	overdueParam := c.Query("overdue")
	notDueParam := c.Query("notDue")

	var schedules []models.TwSchedule
	query := h.DB.
		Table("tw_schedules").
		Select("DISTINCT tw_schedules.*").
		Joins("JOIN tw_schedule_participants ON tw_schedule_participants.schedule_id = tw_schedules.id").
		Joins("JOIN tw_workspaces ON tw_workspaces.id = tw_schedules.workspace_id")

	now := time.Now()
	currentDate := now.Format("2006-01-02")
	fmt.Println("Current date:", currentDate)
	// Apply filters
	if search != "" {
		query = query.Where("tw_schedules.title LIKE ?", "%"+search+"%")
	}
	if dueParam == "day" {
		query = query.Where("DATE(tw_schedules.start_time) = ?", currentDate)
	} else if dueParam == "week" {
		weekEnd := now.AddDate(0, 0, 6).Format("2006-01-02")
		query = query.Where("DATE(tw_schedules.start_time) >= ? AND DATE(tw_schedules.start_time) <= ?", currentDate, weekEnd)
	} else if dueParam == "month" {
		monthEnd := now.AddDate(0, 0, 29).Format("2006-01-02")
		query = query.Where("DATE(tw_schedules.start_time) >= ? AND DATE(tw_schedules.start_time) <= ?", currentDate, monthEnd)
	}
	if dueCompleteParam == "true" {
		query = query.Where("tw_schedules.status = 'done'")
	}
	if overdueParam == "true" {
		query = query.Where("DATE(tw_schedules.start_time) < CURRENT_DATE AND tw_schedules.status != 'done'")
	}
	if notDueParam == "true" {
		query = query.Where("tw_schedules.start_time IS NULL")
	}

	// Filter by member emails if provided
	if membersParam != "" {
		emails := strings.Split(membersParam, ",")
		query = query.Joins("JOIN tw_workspace_users ON tw_workspace_users.id = tw_schedule_participants.workspace_user_id ").
			Joins("JOIN tw_user_emails ON tw_user_emails.id = tw_workspace_users.user_email_id").
			Where("tw_schedule_participants.invitation_status ='joined' AND tw_schedule_participants.deleted_at IS NULL").
			Where("tw_workspace_users.deleted_at IS NULL AND tw_workspace_users.status = 'joined' AND tw_workspace_users.is_active=true AND tw_workspace_users.is_verified = true").
			Where("tw_user_emails.deleted_at IS NULL ").
			Where("tw_user_emails.email IN (?)", emails)
	}
	query = query.
		Where("tw_schedules.board_column_id = ? AND tw_schedules.workspace_id = ? AND tw_schedules.is_deleted = false AND tw_workspaces.deleted_at IS NULL", boardColumnID, workspaceID)

	if result := query.
		Order("tw_schedules.position").
		Find(&schedules); result.Error != nil {
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
