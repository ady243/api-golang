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
	if err := db.AutoMigrate(&models.Users{}, &models.Matches{}, &models.MatchPlayers{}, &models.Analyst{}); err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Initialize services and controllers
	imageService := services.NewImageService("./uploads")
	emailService := services.NewEmailService()
	chatService := services.NewChatService(db, redisClient)
	matchService := services.NewMatchService(db, chatService, redisClient)
	authService := services.NewAuthService(db, imageService, emailService)
	matchPlayersService := services.NewMatchPlayersService(db)
	analystService := services.NewAnalystService(db)
	openAIService := services.NewOpenAIService()

	authController := controllers.NewAuthController(authService, imageService, matchService)
	matchController := controllers.NewMatchController(matchService, authService, db, chatService)
	matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db)
	chatController := controllers.NewChatController(chatService)
	openAiController := controllers.NewOpenAiController(openAIService, matchPlayersService)
	analystController := controllers.NewAnalystController(analystService, authService, db)

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
	routes.SetupRoutesAnalyst(app, analystController)

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Server started on port %s", port)

	// Background task for updating match statuses
	go func() {
		for {
			if err := matchService.UpdateMatchStatuses(); err != nil {
				log.Printf("Error updating match statuses: %v", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	log.Fatal(app.Listen(":" + port))
}
