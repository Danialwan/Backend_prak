package routes

import (
	"project-crud/controllers"
	middlewares "project-crud/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RouteApp(app *fiber.App) {
	api := app.Group("/api")
	// Public routes
	api.Post("/login", controllers.Login)
	api.Get("/", controllers.HomeFunc)
	// api.Get("/all", controllers.GetUsers)
	api.Post("/create", controllers.CreateUser)

	// Protected routes dengan middleware
	users := api.Group("/users", middlewares.AuthMiddleware) // gunakan api.Group untuk nested routes
	users.Get("/", controllers.GetUsers)
	users.Get("/:id", controllers.GetUserByID)
	// users.Post("/", controllers.CreateUser)
	users.Put("/:id", controllers.UpdateUserByID)
	users.Put("/:id/changePass", controllers.UpdatePassword)
	users.Get("/account", controllers.GetUserData) // Atau rute lainnya yang sesuai

	// Contoh route dengan middleware checkRole
	users.Get("/admin", middlewares.CheckRole(1), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome Admin!"})
	})

	// Contoh route dengan middleware checkJenis_user
	users.Get("/mahasiswa", middlewares.CheckJenisUser(1), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome Mahasiswa!"})
	})

	// Route Admin
	// admin := api.Group("/admin", middlewares.AuthMiddleware, middlewares.CheckRole(1)) // Hanya admin
	admin := api.Group("/admin", middlewares.AuthMiddleware)
	admin.Post("/modul", controllers.CreateModul)
	admin.Put("/modul/:id", controllers.UpdateModul)
	admin.Delete("/modul/:id", controllers.DeleteModul)
	admin.Post("/jenis_user/", controllers.AddJenisUser)
	admin.Put("/jenis_user/:id_jenis_user", controllers.UpdateJenisUserModul)
	admin.Delete("/jenis_user/:id_jenis_user", controllers.DeleteJenisUserModul)
	admin.Put("/pindah_user/:user_id", controllers.PindahJenisUser)
	admin.Post("/user/:user_id/modul", controllers.AddModulToUser)

}
