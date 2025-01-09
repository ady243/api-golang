package controllers

import (
	"log"

	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Status string `json:"status"`
}

type ChatController struct {
	ChatService         *services.ChatService
	NotificationService *services.NotificationService
}

func NewChatController(chatService *services.ChatService, notificationService *services.NotificationService) *ChatController {
	return &ChatController{
		ChatService:         chatService,
		NotificationService: notificationService,
	}
}

// @Summary SendMessage
// @Description Send a message to a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param match_id body string true "Match ID"
// @Param user_id body string true "User ID"
// @Param message body string true "Message"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/chat/send [post]
func (ctrl *ChatController) SendMessage(c *fiber.Ctx) error {
	var req struct {
		MatchID  string `json:"match_id" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		Message  string `json:"message" binding:"required"`
		FCMToken string `json:"fcm_token" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request"})
	}

	if err := ctrl.ChatService.AddMessage(req.MatchID, req.UserID, req.Message); err != nil {
		log.Printf("Error saving message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Could not save message"})
	}

	// Récupérer les tokens FCM des participants
	participants, err := ctrl.ChatService.GetParticipants(req.MatchID)
	if err != nil {
		log.Printf("Error fetching participants: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Could not fetch participants"})
	}

	// Envoyer une notification push aux participants
	for _, participant := range participants {
		if participant.ID != req.UserID {
			err := ctrl.NotificationService.SendPushNotification(
				participant.FCMToken,
				"TeamUp",
				"Vous avez reçu un nouveau message dans votre match ",
			)
			if err != nil {
				log.Printf("Failed to send push notification: %v", err)
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Status: "Message sent"})
}

// @Summary GetMessages
// @Description Get messages from a chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param matchID path string true "Match ID"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/chat/{matchID} [get]
func (ctrl *ChatController) GetMessages(c *fiber.Ctx) error {
	matchID := c.Params("matchID")
	limit := 10
	messages, err := ctrl.ChatService.GetMessages(matchID, limit)
	if err != nil {
		log.Printf("Error retrieving messages: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Could not retrieve messages"})
	}

	return c.JSON(messages)
}
