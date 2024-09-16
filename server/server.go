package server

//make just a little server 
import (
    "log" // Ajoutez cette ligne
    "github.com/ady243/teamup/models"
    "github.com/ady243/teamup/storage"
	"github.com/ady243/teamup/services"
	"github.com/ady243/teamup/controllers"
	"github.com/ady243/teamup/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/limiter"
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

	) 
	if err != nil {
		log.Fatal("could not migrate the database", err)
	}


	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)



	// Initialiser l'application Fiber
	app := fiber.New()

	app.Use(helmet.New())

	app.Use(cors.New())

	app.Use(limiter.New())

	routes.SetupRoutesAuth(app, authController)

	// Démarrer le serveur
	log.Fatal(app.Listen(":3003"))

}