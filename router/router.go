package router

import (
	"auth-api-go/controller"
	"auth-api-go/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api/auth", logger.New())

	api.Post("/login", controller.Login) // Authenticate User

	// User Group Route
	user := api.Group("/user")
	user.Get("/:id", controller.GetUser)                               // Get user
	user.Post("/", controller.CreateUser)                              // Create user
	user.Patch("/:id", middleware.Protected(), controller.UpdateUser)  // Update user
	user.Delete("/:id", middleware.Protected(), controller.DeleteUser) // Delete user
}
