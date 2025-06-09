// routes/routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"

	"booksrendezvous-backend/controllers" // Replace "your-module-name" with the actual module name
	"booksrendezvous-backend/middleware"
)

// SetUpRoutes sets up all the routes for the application
func SetUpRoutes(app *fiber.App) {

	// Test route to verify application setup
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", middleware.Protected(), controllers.User)
	//app.Post("/api/simplelogin", controllers.SimpleLogin)
	app.Post("/logout", middleware.Protected(), controllers.Logout)
	app.Post("/api/passwordchange", middleware.Protected(), controllers.PasswordChange)

	// Public user routes
	app.Post("/api/publicuser", controllers.GetPublicUser)
	app.Post("/api/changepublicvisibility", middleware.Protected(), controllers.ChangePublicVisibility)
	app.Get("/api/getpublicvisibility", middleware.Protected(), controllers.GetPublicVisibility)

	app.Get("/api/books", middleware.Protected(), controllers.GetBooks)
	app.Post("/api/addbook", middleware.Protected(), controllers.AddBook)
	app.Delete("/api/books/:id", middleware.Protected(), controllers.DeleteBook)
	app.Put("/api/books/:id", middleware.Protected(), controllers.UpdateBook)

	// stats
	app.Get("/api/stats", middleware.Protected(), controllers.GetStats)

	// achievements
	app.Get("/api/achievements", middleware.Protected(), controllers.GetAchievements)

	// Password reset routes
	app.Post("/api/forgetpassword", controllers.ForgetPassword)
	app.Post("/api/verify-reset-token", controllers.VerifyResetToken)
	app.Post("/api/reset-password", controllers.ResetPassword)

	// google oauth
	//app.Get("/auth/google/callback", controllers.GoogleCallback)

	// linkedin oauth
	//app.Get("/auth/linkedin/callback", controllers.LinkedInCallback)

	// github oauth
	//app.Get("/auth/github/callback", controllers.GitHubCallback)

	// facebook oauth
	//app.Get("/auth/facebook/callback", controllers.FacebookCallback)
}
