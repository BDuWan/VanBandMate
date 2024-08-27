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
	app.Get("/verify", controllers.VerifyHandler)
	app.Get("/forgot-password", controllers.GetForgotPasswordPage)
	app.Post("/login", controllers.PostLogin)
	app.Post("/signup", controllers.PostSignup)
	app.Post("/check-password", controllers.CheckPasswordHandler)
	app.Post("/send-otp", controllers.SendOtp)
	app.Put("/update-password", controllers.PutUpdatePassword)

	//app.Put("/info/update", controllers.PutUpdateUserInformation)
	info := app.Group("/info", IsAuthenticated, CheckSession)
	info.Get("", controllers.GetProfile)
	info.Get("/profile", controllers.GetProfile)
	info.Get("/profile/:id", controllers.GetProfileID)
	info.Get("/userinfo", controllers.GetUserInfo)
	info.Get("/change-password", controllers.GetChangePasswordPage)
	info.Put("/change-password", controllers.PutChangePassword)
	info.Put("/userinfo", controllers.PutUserInfo)
	info.Post("/verify-email", controllers.PostVerifyEmail)

	mngAccount := app.Group("/mng-account", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngAccount)
	mngAccount.Get("", controllers.GetRolePage)

	mngRole := app.Group("/mng-role", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngRole)
	mngRole.Get("", controllers.GetRolePage)
	mngRole.Get("/api", controllers.APIGetRole)
	mngRole.Get("/api/:id", controllers.APIGetRoleID)
	mngRole.Put("/api/update/:id", controllers.APIPutUpdateRoleID)
	mngRole.Post("/api/create", controllers.APIPostCreateRole)
	mngRole.Delete("/api/delete/:id", controllers.APIDeleteRoleID)

	mngContract := app.Group("/mng-contract", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngContract)
	mngContract.Get("", controllers.GetRolePage)

	mngHiringNews := app.Group("/mng-hiring-news", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngHiringNews)
	mngHiringNews.Get("", controllers.GetRolePage)

	mngDashboard := app.Group("/dashboard", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionDashboard)
	mngDashboard.Get("", controllers.GetRolePage)

	contract := app.Group("/contract", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMyContract)
	contract.Get("", controllers.GetMyContractPage)
	contract.Get("/api", controllers.APIGetContract)
	contract.Get("/api/detail/:id", controllers.APIGetContractDetailID)
	contract.Post("/api/filter", controllers.APIPostContractFilter)
	contract.Post("/request-delete", controllers.PostContractRequestDelete)
	contract.Post("/confirm-delete", controllers.PostContractConfirmDelete)
	contract.Post("/cancel-delete", controllers.PostContractCancelDelete)

	hiring := app.Group("/hiring", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionHiring)
	hiring.Get("", controllers.GetHiringPage)
	hiring.Get("/api", controllers.APIGetHiring)
	hiring.Get("/api/:id", controllers.APIGetHiringID)
	hiring.Get("/api/detail/:id", controllers.APIGetHiringDetailID)
	hiring.Get("/api/list-apply/:id", controllers.APIGetHiringListApply)
	hiring.Post("/api/filter", controllers.APIPostHiringFilter)
	hiring.Post("/api/create", controllers.APIPostHiringCreate)
	hiring.Post("/api/save-apply", controllers.APIPostSaveApply)
	hiring.Put("/api/edit/:id", controllers.APIPutHiringUpdate)

	findMusicPlayer := app.Group("/find-music-player", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionFindMusicPlayer)
	findMusicPlayer.Get("", controllers.GetRolePage)

	news := app.Group("/news", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionNews)
	news.Get("", controllers.GetNewsPage)
	news.Get("/api", controllers.APIGetNews)
	news.Get("/api/detail/:id", controllers.APIGetHiringDetailID)
	news.Post("/api/filter", controllers.APIPostNewsFilter)
	news.Post("/apply", controllers.PostNewsApply)
	news.Post("/cancel-apply", controllers.PostNewsCancelApply)

	invitation := app.Group("/invitation", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionInvitation)
	invitation.Get("", controllers.GetRolePage)

	errors := app.Group("/errors", IsAuthenticated, CheckSession)
	errors.Get("/403", controllers.GetError403)
	errors.Get("/401", controllers.GetError401)
	errors.Get("/404", controllers.GetError404)
}
