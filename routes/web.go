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
	home.Get("", controllers.GetHome)
	home.Get("/profile/:id", controllers.GetProfileID)
	home.Get("/userinfo", controllers.GetUserInfo)
	home.Put("/userinfo", controllers.PutUserInfo)

	mngRoles := app.Group("/mng-role", IsAuthenticated, CheckSession, CheckPermissionMngRole)
	mngRoles.Get("", controllers.GetRolePage)
	mngRoles.Get("/api", controllers.APIGetRole)
	mngRoles.Get("/api/:id", controllers.APIGetRoleID)
	mngRoles.Post("/api/create", controllers.APIPostCreateRole)

	errors := app.Group("/errors", IsAuthenticated, CheckSession)
	errors.Get("/403", controllers.GetError403)
	errors.Get("/401", controllers.GetError401)
	errors.Get("/404", controllers.GetError404)
}
