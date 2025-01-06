package controllers

import (
	"log"

	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	NotificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
	}
}

func (nc *NotificationController) SendPushNotification(c *fiber.Ctx) error {
	var request struct {
		Token string `json:"token"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := nc.NotificationService.SendPushNotification(request.Token, request.Title, request.Body); err != nil {
		log.Printf("Error sending notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification sent successfully",
	})
}
