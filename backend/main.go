// main.go
package main

import (
	"booksrendezvous-backend/controllers"
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/routes"
	"booksrendezvous-backend/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load environment variables
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger with the log level from environment variables
	err = utils.Initialize(config.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.Sync()

	sugar := utils.SugaredLogger

	sugar.Infow("Starting application with config",
		"DBHost", config.DBHost,
		"DBPort", config.DBPort,
		"FrontendURL", config.FrontendURL,
		"ServerPort", config.ServerPort,
	)

	// Connect to the database
	db, err := database.ConnectDB(config, sugar)
	if err != nil {
		sugar.Fatalw("Failed to connect to database",
			"error", err,
		)
	}
	sugar.Info("Successfully connected to database")

	// Initialize Fiber app
	app := fiber.New()

	// Adding CORS middleware with specific origin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.FrontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Accept, credentials",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", config.FrontendURL)
		c.Set("Access-Control-Allow-Credentials", "true")
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Process the next middleware/handler
		err := c.Next()

		// Log the request details
		log.Printf("[%s] %s %s - %s %s",
			time.Now().Format("2006-01-02 15:04:05"), // Timestamp
			c.Method(),                               // HTTP Method (GET, POST, etc.)
			c.OriginalURL(),                          // Requested URL
			c.IP(),                                   // Client IP
			time.Since(startTime),                    // Processing time
		)

		return err
	})

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusOK) // Let the OPTIONS request pass through
		}
		return c.Next() // Proceed to the next handler
	})

	// Setup routes
	routes.SetUpRoutes(app)

	// add base success

	err = controllers.SeedBaseAchievements(db, "./data/succes.json")
	if err != nil {
		panic(err)
	}

	sugar.Infof("Server starting on port%s", config.ServerPort)
	if err := app.Listen(config.ServerPort); err != nil {
		sugar.Fatalw("Server failed to start",
			"error", err,
		)
	}
}
