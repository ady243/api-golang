package server

import (
	"log"
	"os"
	"time"

	_ "github.com/ady243/teamup/docs"
	"github.com/ady243/teamup/internal"
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
	if err := db.AutoMigrate(&models.Users{}, &models.Matches{}, &models.MatchPlayers{}, &models.FriendRequest{}, &models.Message{}); err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Initialize Firebase
	_, err = internal.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize services and controllers
	imageService := services.NewImageService("./uploads")
	emailService := services.NewEmailService()
	matchService := services.NewMatchService(db, services.NewChatService(db, redisClient), redisClient)
	authService := services.NewAuthService(db, imageService, emailService)
	webSocketService := services.NewWebSocketService()
	notificationService, err := services.NewNotificationService()
	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}

	openAIService := services.NewOpenAIService()
	friendChatService := services.NewFriendChatService(db, webSocketService, notificationService)

	matchPlayersService := services.NewMatchPlayersService(db)
	friendService := services.NewFriendService(db, authService, webSocketService)
	friendController := controllers.NewFriendController(friendService, notificationService)
	chatService := services.NewChatService(db, redisClient)
	matchController := controllers.NewMatchController(matchService, authService, db, chatService, redisClient, matchPlayersService)
	matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db)
	chatController := controllers.NewChatController(chatService, notificationService)
	openAiController := controllers.NewOpenAiController(openAIService, matchPlayersService)
	authController := controllers.NewAuthController(authService, imageService, matchService)
	friendChatController := controllers.NewFriendChatController(friendChatService, friendService, notificationService)
	notificationController := controllers.NewNotificationController(notificationService)

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
	routes.SetupFriendRoutes(app, friendController)
	routes.SetupRoutesFriendMessage(app, friendChatController)
	routes.SetupNotificationRoutes(app, notificationController)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		webSocketService.HandleWebSocket(c)
	}))

	// Start WebSocket broadcast
	go webSocketService.StartBroadcast()

	// Start listening for notifications
	// go notificationService.ListenForNotifications()

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Server started on port %s", port)

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
