package server

import (
	"log"
	"os"

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
	if err := db.AutoMigrate(&models.Users{}, &models.Matches{}, &models.MatchPlayers{}); err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	// Optional: Create fake users for testing
	// users := models.GenerateFakeUsers(10)
	// for _, user := range users {
	//     db.Create(&user)
	// }

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// Initialize services and controllers
	imageService := services.NewImageService("./uploads")
	authService := services.NewAuthService(db, imageService, services.NewEmailService())
	authController := controllers.NewAuthController(authService, imageService)
	matchService := services.NewMatchService(db)
	openAIService := services.NewOpenAIService()
	chatService := services.NewChatService(db, redisClient)
	matchController := controllers.NewMatchController(matchService, authService, db, chatService)
	matchPlayersService := services.NewMatchPlayersService(db)
	matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db)
	chatController := controllers.NewChatController(chatService)
	openAiController := controllers.NewOpenAiController(openAIService, matchPlayersService)

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
			Expiration: 30 * 1000,
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

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
