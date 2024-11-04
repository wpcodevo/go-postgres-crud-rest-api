package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/go-postgres-crud-rest-api/initializers"
	"github.com/wpcodevo/go-postgres-crud-rest-api/models"
	"gorm.io/gorm"
)

func CreateFeedbackHandler(c *fiber.Ctx) error {
	var payload *models.CreateFeedbackSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	now := time.Now()
	newFeedback := models.Feedback{
		Name:      payload.Name,
		Email:     payload.Email,
		Feedback:  payload.Feedback,
		Rating:    payload.Rating,
		Status:    payload.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newFeedback)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Feedback already exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": newFeedback}})
}

func FindFeedbacksHandler(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var feedbacks []models.Feedback
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&feedbacks)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(feedbacks), "feedbacks": feedbacks})
}

func UpdateFeedbackHandler(c *fiber.Ctx) error {
	feedbackId := c.Params("feedbackId")

	var payload *models.UpdateFeedbackSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var feedback models.Feedback
	result := initializers.DB.First(&feedback, "id = ?", feedbackId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No feedback with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.Feedback != "" {
		updates["feedback"] = payload.Feedback
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}

	if payload.Rating != nil {
		updates["rating"] = payload.Rating
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&feedback).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"feedback": feedback}})
}

func FindFeedbackByIdHandler(c *fiber.Ctx) error {
	feedbackId := c.Params("feedbackId")

	var feedback models.Feedback
	result := initializers.DB.First(&feedback, "id = ?", feedbackId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No feedback with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"feedback": feedback}})
}

func DeleteFeedbackHandler(c *fiber.Ctx) error {
	feedbackId := c.Params("feedbackId")

	result := initializers.DB.Delete(&models.Feedback{}, "id = ?", feedbackId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
