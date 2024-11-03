package controllers

import (
    "log"

    "github.com/ady243/teamup/internal/services"
    "github.com/gofiber/fiber/v2"
)

type ChatController struct {
    ChatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
    return &ChatController{
        ChatService: chatService,
    }
}

func (ctrl *ChatController) SendMessage(c *fiber.Ctx) error {
    var req struct {
        MatchID string `json:"match_id"`
        UserID  string `json:"user_id"`
        Message string `json:"message"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    if err := ctrl.ChatService.AddMessage(req.MatchID, req.UserID, req.Message); err != nil {
        log.Printf("Error saving message: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save message"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Message sent"})
}

func (ctrl *ChatController) GetMessages(c *fiber.Ctx) error {
    matchID := c.Params("matchID")
    limit := 10
    messages, err := ctrl.ChatService.GetMessages(matchID, limit)
    if err != nil {
        log.Printf("Error retrieving messages: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve messages"})
    }

    return c.JSON(messages)
}