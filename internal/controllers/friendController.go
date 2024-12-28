package controllers

import (
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FriendController struct {
	friendService       *services.FriendService
	notificationService *services.NotificationService
}

func NewFriendController(friendService *services.FriendService, notificationService *services.NotificationService) *FriendController {
	return &FriendController{friendService: friendService, notificationService: notificationService}
}

func (c *FriendController) SendFriendRequest(ctx *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.friendService.SendFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification
	if err := c.notificationService.SendNotification(request.ReceiverId, "New Friend Request", "You have a new friend request!"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Friend request sent"})
}

func (c *FriendController) AcceptFriendRequest(ctx *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.friendService.AcceptFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification
	if err := c.notificationService.SendNotification(request.SenderId, "Friend Request Accepted", "Your friend request has been accepted!"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Friend request accepted"})
}

func (c *FriendController) GetFriends(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	friends, err := c.friendService.GetFriends(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(friends)
}

func (c *FriendController) SearchUsersByUsername(ctx *fiber.Ctx) error {
	username := ctx.Query("username")

	users, err := c.friendService.SearchUsersByUsername(username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(users)
}
