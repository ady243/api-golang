package controllers

import (
	"log"

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
		log.Printf("Failed to parse request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("Sending friend request from %s to %s", request.SenderId, request.ReceiverId)
	if err := c.friendService.SendFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		log.Printf("Failed to send friend request: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification
	if err := c.notificationService.SendNotification(request.ReceiverId, "New Friend Request", "You have a new friend request!"); err != nil {
		log.Printf("Failed to send notification: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Friend request sent successfully from %s to %s", request.SenderId, request.ReceiverId)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Friend request sent"})
}

func (c *FriendController) AcceptFriendRequest(ctx *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("Accepting friend request from %s to %s", request.SenderId, request.ReceiverId)
	if err := c.friendService.AcceptFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		log.Printf("Failed to accept friend request: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification
	if err := c.notificationService.SendNotification(request.SenderId, "Friend Request Accepted", "Your friend request has been accepted!"); err != nil {
		log.Printf("Failed to send notification: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Friend request accepted successfully from %s to %s", request.SenderId, request.ReceiverId)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Friend request accepted"})
}

func (c *FriendController) GetFriends(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	friends, err := c.friendService.GetFriends(userID)
	if err != nil {
		log.Printf("Failed to get friends for user %s: %v", userID, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(friends)
}

func (c *FriendController) GetFriendRequests(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	requests, err := c.friendService.GetFriendRequests(userID)
	if err != nil {
		log.Printf("Failed to get friend requests for user %s: %v", userID, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(requests)
}

func (c *FriendController) SearchUsersByUsername(ctx *fiber.Ctx) error {
	username := ctx.Query("username")
	users, err := c.friendService.SearchUsersByUsername(username)
	if err != nil {
		log.Printf("Failed to search users by username %s: %v", username, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(users)
}
