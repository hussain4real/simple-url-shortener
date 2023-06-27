package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hussain4real/simple-url-shortener/controllers"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	//grouping
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//routes

	v1.Get("/home", controllers.Home)
	v1.Get("/user", controllers.GetUser)
	v1.Post("/register", controllers.RegisterUser)
	v1.Post("/login", controllers.LoginUser)
	v1.Post("/logout", controllers.LogoutUser)
	v1.Put("/users/:id", controllers.UpdateUser)
	v1.Delete("/users/:id", controllers.DeleteUser)

	//shortly routes
	v1.Get("/r/:redirect", controllers.Redirect)
	v1.Get("/shortly", controllers.GetAllShortlies)
	v1.Get("/shortly/:id", controllers.GetShortly)
	v1.Post("/shortly", controllers.CreateShortly)

}
