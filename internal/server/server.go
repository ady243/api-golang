package server

import (
	"log"
	"os"
	"time"

	_ "github.com/ady243/teamup/docs"
	"github.com/ady243/teamup/internal/controllers"
	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/routes"
	"github.com/ady243/teamup/internal/services"
	"github.com/ady243/teamup/storage"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan services.Notification)

func Run() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found", err)
	}

	// Database connection
	db, err := storage.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Table migration
	if err := db.AutoMigrate(&models.Users{}, &models.Matches{}, &models.MatchPlayers{}, &models.Friendship{}, &models.FriendRequest{}); err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Initialize services and controllers
	imageService := services.NewImageService("./uploads")
	emailService := services.NewEmailService()
	matchService := services.NewMatchService(db, services.NewChatService(db, redisClient), redisClient)
	authService := services.NewAuthService(db, imageService, emailService)
	friendService := services.NewFriendService(db, redisClient, authService, nil)
	notificationService := services.NewNotificationService(redisClient, broadcast)
	authController := controllers.NewAuthController(authService, imageService, matchService)
	friendController := controllers.NewFriendController(friendService, notificationService)
	matchService = services.NewMatchService(db, services.NewChatService(db, redisClient), redisClient)
	openAIService := services.NewOpenAIService()
	chatService := services.NewChatService(db, redisClient)
	webSocketService := services.NewWebSocketService()
	matchController := controllers.NewMatchController(matchService, authService, db, chatService, redisClient)
	matchPlayersService := services.NewMatchPlayersService(db)
	matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db)
	chatController := controllers.NewChatController(chatService)
	openAiController := controllers.NewOpenAiController(openAIService, matchPlayersService)
	webSocketController := controllers.NewWebSocketController(webSocketService)

	// Configure Fiber app
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Apply rate limiter middleware to all routes except Swagger
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/swagger/*" {
			return c.Next()
		}
		return limiter.New(limiter.Config{
			Max:        10,
			Expiration: 30 * time.Second,
		})(c)
	})

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to TeamUp API!")
	})
	routes.SetupRoutesAuth(app, authController)
	routes.SetupRoutesMatches(app, matchController)
	routes.SetupRoutesMatchePlayers(app, matchPlayersController)
	routes.SetupChatRoutes(app, chatController)
	routes.SetupOpenAiRoutes(app, openAiController)
	routes.SetupRoutesWebSocket(app, webSocketController)
	routes.SetupRoutesFriend(app, friendController)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clients[c] = true
		defer func() {
			delete(clients, c)
			c.Close()
		}()

		for {
			var msg services.Notification
			if err := c.ReadJSON(&msg); err != nil {
				log.Printf("error: %v", err)
				break
			}
			broadcast <- msg
		}
	}))

	go handleMessages()

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Server started on port %s", port)
	go webSocketService.StartBroadcast()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := matchService.UpdateMatchStatuses(); err != nil {
				log.Printf("Erreur lors de la mise à jour des statuts des matchs : %v", err)
			}
		}
	}()

	log.Fatal(app.Listen(":" + port))
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("websocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
