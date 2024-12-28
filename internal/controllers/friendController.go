package controllers

import (
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FriendController struct {
	FriendService       *services.FriendService
	NotificationService *services.NotificationService
}

func NewFriendController(friendService *services.FriendService, notificationService *services.NotificationService) *FriendController {
	return &FriendController{
		FriendService:       friendService,
		NotificationService: notificationService,
	}
}

func (fc *FriendController) SendFriendRequest(c *fiber.Ctx) error {
	type request struct {
		SenderID   string `json:"senderId"`
		ReceiverID string `json:"receiverId"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.SendFriendRequest(req.SenderID, req.ReceiverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	title := "Nouvelle demande d'ami"
	message := "Vous avez reçu une nouvelle demande d'ami"
	if err := fc.NotificationService.SendNotification(req.ReceiverID, title, message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "friend request sent",
	})
}

func (fc *FriendController) AcceptFriendRequest(c *fiber.Ctx) error {
	type request struct {
		SenderID   string `json:"senderId"`
		ReceiverID string `json:"receiverId"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.AcceptFriendRequest(req.SenderID, req.ReceiverID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	title := "Demande d'ami acceptée"
	message := "Votre demande d'ami a été acceptée"
	if err := fc.NotificationService.SendNotification(req.SenderID, title, message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "friend request accepted",
	})
}

func (fc *FriendController) GetFriends(c *fiber.Ctx) error {
	userID := c.Params("userID")
	friends, err := fc.FriendService.GetFriends(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(friends)
}

func (fc *FriendController) SearchUsersByUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	users, err := fc.FriendService.SearchUsersByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(users)
}
