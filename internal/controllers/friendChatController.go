package controllers

import (
	"log"

	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FriendChatController struct {
	FriendChatService *services.FriendChatService
	FriendService     *services.FriendService
}

func NewFriendChatController(friendChatService *services.FriendChatService, friendService *services.FriendService) *FriendChatController {
	return &FriendChatController{
		FriendChatService: friendChatService,
		FriendService:     friendService,
	}
}

func (cc *FriendChatController) SendMessage(c *fiber.Ctx) error {
	var request struct {
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Vérifiez que sender_id et receiver_id ne sont pas vides
	if request.SenderID == "" || request.ReceiverID == "" {
		log.Printf("SenderID or ReceiverID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be empty",
		})
	}

	// Vérifiez que sender_id et receiver_id ne sont pas identiques
	if request.SenderID == request.ReceiverID {
		log.Printf("SenderID and ReceiverID are the same")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be the same",
		})
	}

	// Vérifiez que les utilisateurs sont amis
	areFriends, err := cc.FriendService.AreFriends(request.SenderID, request.ReceiverID)
	if err != nil || !areFriends {
		log.Printf("Users are not friends or error occurred: %v", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "users are not friends",
		})
	}

	// Envoi du message
	if err := cc.FriendChatService.SendMessage(request.SenderID, request.ReceiverID, request.Content); err != nil {
		log.Printf("Error sending message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "message sent successfully",
	})
}

func (cc *FriendChatController) GetMessages(c *fiber.Ctx) error {
	senderID := c.Params("senderID")
	receiverID := c.Params("receiverID")

	areFriends, err := cc.FriendService.AreFriends(senderID, receiverID)
	if err != nil || !areFriends {
		log.Printf("Users are not friends or error occurred: %v", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "users are not friends",
		})
	}

	messages, err := cc.FriendChatService.GetMessages(senderID, receiverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(messages)
}
