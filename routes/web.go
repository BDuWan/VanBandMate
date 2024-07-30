package routes

import (
	"lms/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == "POST" {
			method := c.FormValue("_method")

			switch method {
			case "DELETE":
				c.Method("DELETE")
			case "PUT":
				c.Method("PUT")
			}
		}
		return c.Next()
	})

	app.Get("/", controllers.GetLogin)
	app.Get("/login", controllers.GetLogin)
	app.Get("/signup", controllers.GetSignup)
	app.Get("/logout", controllers.GetLogout)
	app.Post("/login", controllers.PostLogin)
	app.Post("/forgot-password", controllers.CreateNewPass)
	app.Post("/signup", controllers.PostSignup)

	//app.Put("/home/update", controllers.PutUpdateUserInformation)
	home := app.Group("/home", IsAuthenticated, CheckSession)
	//home.Get("/information", controllers.GetInformation)
	//home.Get("/changePassword", controllers.GetChangePassword)
	home.Get("", controllers.GetHome)

	accounts := app.Group("/accounts", IsAuthenticated, CheckSession, CheckVerify)

	accUsers := accounts.Group("/users", CheckPerAccUsers)
	accUsers.Get("", controllers.GetAccUsers)
	accUsers.Get("/create", controllers.GetAccCreateUser)
	accUsers.Get("/:id", controllers.GetAccUserID)
	accUsers.Post("", controllers.PostAccCreateUser)
	accUsers.Post("/api/admin", controllers.APIPostAccUserAdmin)
	accUsers.Post("/api/sales", controllers.APIPostAccUserSales)
	accUsers.Post("/api/business", controllers.APIPostAccUserBusiness)
	accUsers.Post("/api/instructors", controllers.APIPostAccUserInstructors)
	accUsers.Post("/api/students", controllers.APIPostAccUserStudents)
	accUsers.Delete("/:id", controllers.DeleteAccUserID)
	accUsers.Put("/:id", controllers.UpdateAccUserID)
	accUsers.Put("/state/:id", controllers.UpdateAccStateUser)

	accRoles := accounts.Group("/roles", CheckPerAccRoles)
	accRoles.Get("", controllers.GetAccRoles)
	accRoles.Get("/create", controllers.GetAccCreateRole)
	accRoles.Get("/:id", controllers.GetAccRoleID)
	accRoles.Post("", controllers.PostAccCreateRole)
	accRoles.Post("/api", controllers.APIPostAccRoles)
	accRoles.Delete("/:id", controllers.DeleteAccRoleID)
	accRoles.Put("/:id", controllers.UpdateAccRoleID)

	errors := app.Group("/errors", IsAuthenticated, CheckSession)
	errors.Get("/403", controllers.GetError403)
	errors.Get("/401", controllers.GetError401)
	errors.Get("/404", controllers.GetError404)
}
