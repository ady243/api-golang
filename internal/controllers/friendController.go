package controllers

import (
	"log"

	"github.com/ady243/teamup/internal/models"
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
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if request.SenderId == "" || request.ReceiverId == "" {
		log.Printf("SenderId or ReceiverId is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be empty",
		})
	}

	if request.SenderId == request.ReceiverId {
		log.Printf("SenderId and ReceiverId are the same")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sender_id and receiver_id cannot be the same",
		})
	}

	if err := fc.FriendService.SendFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		log.Printf("Error sending friend request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var receiver models.Users
	if err := fc.FriendService.DB.Where("id = ?", request.ReceiverId).First(&receiver).Error; err != nil {
		log.Printf("Failed to fetch receiver: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Receiver not found",
		})
	}

	// Envoyer une notification push au destinataire
	err := fc.NotificationService.SendPushNotification(
		receiver.FCMToken,
		"TeamUp rélation",
		"Vous avez reçu une nouvelle demande d'ami de "+receiver.Username,
	)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "friend request sent successfully",
	})
}

func (fc *FriendController) AcceptFriendRequest(c *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.AcceptFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Récupérer le token FCM de l'expéditeur
	var sender models.Users
	if err := fc.FriendService.DB.Where("id = ?", request.SenderId).First(&sender).Error; err != nil {
		log.Printf("Failed to fetch sender: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sender not found",
		})
	}

	// Envoyer une notification push à l'expéditeur
	err := fc.NotificationService.SendPushNotification(
		sender.FCMToken,
		"TeamUp rélation",
		"Votre demande d'ami a été acceptée par "+sender.Username,
	)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "friend request accepted successfully",
	})
}

func (fc *FriendController) DeclineFriendRequest(c *fiber.Ctx) error {
	var request struct {
		SenderId   string `json:"sender_id"`
		ReceiverId string `json:"receiver_id"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := fc.FriendService.DeclineFriendRequest(request.SenderId, request.ReceiverId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Récupérer le token FCM de l'expéditeur
	var sender models.Users
	if err := fc.FriendService.DB.Where("id = ?", request.SenderId).First(&sender).Error; err != nil {
		log.Printf("Failed to fetch sender: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sender not found",
		})
	}

	// Envoyer une notification push à l'expéditeur
	err := fc.NotificationService.SendPushNotification(
		sender.FCMToken,
		"TeamUp rélation",
		"Votre demande d'ami a été refusée par "+sender.Username,
	)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "friend request declined successfully",
	})
}

func (fc *FriendController) GetFriendRequests(c *fiber.Ctx) error {
	userID := c.Params("userID")
	friendRequests, err := fc.FriendService.GetFriendRequests(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(friendRequests)
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
