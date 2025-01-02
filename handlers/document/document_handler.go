package document

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/document_dtos"
	"github.com/timewise-team/timewise-models/models"
	"gorm.io/gorm"
	"log"
)

// getDocumentsBySchedule godoc
// @Summary Get documents by schedule
// @Description Get documents by schedule
// @Tags document
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwDocument
// @Router /dbms/v1/document/schedule/{schedule_id} [get]
func (h *DocumentHandler) getDocumentsBySchedule(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")
	if scheduleId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var Documents []models.TwDocument
	if err := h.DB.
		Where("schedule_id = ?", scheduleId).
		Where("deleted_at IS NULL").
		Find(&Documents).Error; err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(Documents)
}

// getDocumentsBySchedule godoc
// @Summary Get documents by schedule
// @Description Get documents by schedule
// @Tags document
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} models.TwDocument
// @Router /dbms/v1/document/schedule_id/{schedule_id} [get]
func (h *DocumentHandler) getDocumentsByScheduleID(c *fiber.Ctx) error {
	scheduleId := c.Params("schedule_id")
	if scheduleId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var documents []document_dtos.TwDocumentResponse
	err := h.DB.Table("tw_documents AS d").
		Select(`
            d.id AS id,
            d.file_name,
            d.file_path,
			d.file_size,
			d.file_type,
			d.is_deleted,
			d.schedule_id,
			d.uploaded_by,
			d.created_at,
			d.updated_at,
			d.uploaded_at,
			d.download_url,
			wu.role,
			wu.status AS status_workspace_user,
			wu.is_verified,
			ue.id as user_id,
			ue.email,
			u.first_name,
			u.last_name,
			u.profile_picture
            
        `).
		Joins("JOIN tw_workspace_users AS wu ON wu.id =d.uploaded_by").
		Joins("JOIN tw_user_emails AS ue ON wu.user_email_id = ue.id").
		Joins("JOIN tw_users AS u ON ue.user_id = u.id").
		Where("d.schedule_id = ?", scheduleId).
		Where("d.deleted_at IS NULL").
		Scan(&documents).Error
	if err != nil {
		log.Println("Error querying document:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách document",
		})
	}

	return c.JSON(documents)
}

// createDocument godoc
// @Summary Create document
// @Description Create document
// @Tags document
// @Accept json
// @Produce json
// @Param document body models.TwDocument true "Document object"
// @Success 200 {object} models.TwDocument
// @Router /dbms/v1/document/upload [post]
func (h *DocumentHandler) createDocument(c *fiber.Ctx) error {
	var document models.TwDocument
	if err := c.BodyParser(&document); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.DB.Create(&document).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(document)
}

// deleteDocument godoc
// @Summary Delete document
// @Description Delete document
// @Tags document
// @Accept json
// @Produce json
// @Param scheduleId query string true "Schedule ID associated with the file"
// @Param fileName query string true "Name of the file to delete"
// @Success 204 "No Content"
// @Router /dbms/v1/document [delete]
func (h *DocumentHandler) deleteDocument(c *fiber.Ctx) error {
	scheduleID := c.Query("scheduleId")
	if scheduleID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	fileName := c.Query("fileName")
	if fileName == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if err := h.DB.Where("schedule_id = ? AND file_name = ?", scheduleID, fileName).Delete(&models.TwDocument{}).Error; err != nil {
		return fmt.Errorf("failed to delete document from database: %v", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// getDocumentsById godoc
// @Summary Get document by ID
// @Description Get document by ID
// @Tags document
// @Accept json
// @Produce json
// @Param document_id path string true "Document ID"
// @Success 200 {object} models.TwDocument
// @Router /dbms/v1/document/{document_id} [get]
func (h *DocumentHandler) getDocumentsById(c *fiber.Ctx) error {
	documentID := c.Params("document_id")
	if documentID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var document models.TwDocument
	if err := h.DB.First(&document, documentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(document)
}
