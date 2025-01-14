package controllers

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/ady243/teamup/internal/services"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
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

    return c.JSON(SuccessResponse{Status: "Message sent successfully"})
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

// @Summary ChatWebSocketHandler
// @Description Handle WebSocket connections for chat
// @Tags Chat
// @Accept json
// @Produce json
// @Param matchID path string true "Match ID"
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/chat/ws/{matchID} [get]
func (ctrl *ChatController) ChatWebSocketHandler(c *websocket.Conn) {
    matchID := c.Params("matchID")
    userID := c.Locals("user_id").(string)

    log.Printf("Handling new WebSocket connection for match: %s", matchID)

    // Vérifier si l'utilisateur est dans le match
    if err := ctrl.ChatService.IsUserInMatch(matchID, userID); err != nil {
        c.Close()
        log.Println("User not in match:", err)
        return
    }

    room := "match:" + matchID
    ctx := context.Background()
    pubsub := ctrl.RedisClient.Subscribe(ctx, room)
    defer pubsub.Close()

    // Goroutine pour écouter les messages de Redis
    go func() {
        for {
            msg, err := pubsub.ReceiveMessage(ctx)
            if err != nil {
                log.Printf("Erreur de réception de message dans Redis : %v", err)
                break
            }
            if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
                log.Printf("Erreur d'envoi de message WebSocket : %v", err)
                break
            }
        }
    }()

    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("Error reading WebSocket message: %v", err)
            break
        }

        // Construire un message avec métadonnées
        fullMessage := fiber.Map{
            "playerId":  userID,
            "message":   string(message),
            "timestamp": time.Now().Format(time.RFC3339),
        }

        // Publier le message au format JSON
        msgJSON, _ := json.Marshal(fullMessage)
        ctrl.RedisClient.Set(ctx, "chat:"+matchID+":"+userID, msgJSON, time.Hour*24*7)
        ctrl.RedisClient.Publish(ctx, room, string(msgJSON))
    }

    log.Printf("WebSocket connection closed for match: %s", matchID)
}