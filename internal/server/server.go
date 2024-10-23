package server

import (
	"log"
	"os"

	"github.com/ady243/teamup/internal/controllers"
	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/routes"
	"github.com/ady243/teamup/internal/services"
	"github.com/ady243/teamup/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Connexion à la base de données
	db, err := storage.NewConnection()
	if err != nil {
		log.Fatal("could not connect to the database", err)
	}

	err = db.AutoMigrate(
		&models.Users{},
		&models.Matches{},
		&models.MatchPlayers{},
	)
	if err != nil {
		log.Fatal("could not migrate the database", err)
	}

	//ici c'est pour generer des faux utilisateurs pour tester
	users := models.GenerateFakeUsers(10)
	for _, user := range users {
		db.Create(&user)
	}

	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)
	matchService := services.NewMatchService(db)
	openAIService := services.NewOpenAIService() 
	matchController := controllers.NewMatchController(matchService, authService, db)
	matchPlayersService := services.NewMatchPlayersService(db)
	matchPlayersController := controllers.NewMatchPlayersController(matchPlayersService, authService, db, openAIService)

	app := fiber.New()

	app.Use(helmet.New(
		
	))

	app.Use(cors.New())

	app.Use(limiter.New())

	routes.SetupRoutesAuth(app, authController)
	routes.SetupRoutesMatches(app, matchController)
	routes.SetupRoutesMatchePlayers(app, matchPlayersController)


	// Démarrer le serveur
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "1234"
	}
	log.Fatal(app.Listen(":" + port))

}
