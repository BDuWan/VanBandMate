package routes

import (
	"vanbandmate/controllers"

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

	mngUser := app.Group("/mng-user", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngUser)
	mngUser.Get("", controllers.GetMngUserPage)
	mngUser.Get("/edit-user/:id", controllers.GetEditUserPage)
	mngUser.Post("/api/filter", controllers.APIPostUserFilter)
	mngUser.Post("/api/create", controllers.APIPostCreateUser)
	mngUser.Put("/api/edit", controllers.APIPutEditUser)
	mngUser.Delete("/api/delete/:id", controllers.APIDeleteUserID)

	mngRole := app.Group("/mng-role", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngRole)
	mngRole.Get("", controllers.GetRolePage)
	mngRole.Get("/api", controllers.APIGetRole)
	mngRole.Get("/api/:id", controllers.APIGetRoleID)
	mngRole.Put("/api/update/:id", controllers.APIPutUpdateRoleID)
	mngRole.Post("/api/create", controllers.APIPostCreateRole)
	mngRole.Delete("/api/delete/:id", controllers.APIDeleteRoleID)

	mngContract := app.Group("/mng-contract", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngContract)
	mngContract.Get("", controllers.GetMngContractPage)
	mngContract.Post("/api/filter", controllers.APIPostMngContractFilter)
	mngContract.Get("/api/detail/:id", controllers.APIGetMngContractDetailID)
	mngContract.Put("/api/restore/:id", controllers.APIPutMngContractRestoreID)
	mngContract.Delete("/api/delete/:id", controllers.APIPutMngContractDeleteID)

	//mngHiringNews := app.Group("/mng-hiring-news", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMngHiringNews)
	//mngHiringNews.Get("", controllers.GetRolePage)

	mngDashboard := app.Group("/dashboard", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionDashboard)
	mngDashboard.Get("", controllers.GetRolePage)

	contract := app.Group("/contract", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionMyContract)
	contract.Get("", controllers.GetMyContractPage)
	//contract.Get("/api", controllers.APIGetContract)
	contract.Get("/api/detail/:id", controllers.APIGetContractDetailID)
	contract.Post("/api/filter", controllers.APIPostContractFilter)
	contract.Post("/request-delete", controllers.PostContractRequestDelete)
	contract.Post("/confirm-delete", controllers.PostContractConfirmDelete)
	contract.Post("/cancel-delete", controllers.PostContractCancelDelete)

	hiring := app.Group("/hiring", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionHiring)
	hiring.Get("", controllers.GetHiringPage)
	hiring.Get("/invite/:id", controllers.GetHiringInvitePage)
	//hiring.Get("/api", controllers.APIGetHiring)
	hiring.Get("/api/:id", controllers.APIGetHiringID)
	hiring.Get("/api/detail/:id", controllers.APIGetHiringDetailID)
	hiring.Get("/api/list-apply/:id", controllers.APIGetHiringListApply)
	hiring.Post("/api/filter", controllers.APIPostHiringFilter)
	hiring.Post("/api/find", controllers.APIPostHiringFind)
	hiring.Post("/api/create", controllers.APIPostHiringCreate)
	hiring.Post("/api/save-apply", controllers.APIPostSaveApply)
	hiring.Post("/invite", controllers.PostHiringInvite)
	hiring.Post("/cancel-invite", controllers.PostHiringCancelInvite)
	hiring.Put("/api/edit/:id", controllers.APIPutHiringUpdate)

	sentInvitation := app.Group("/sent-invitation", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionHiring)
	sentInvitation.Get("", controllers.GetUserInfo)

	news := app.Group("/news", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionNews)
	news.Get("", controllers.GetNewsPage)
	//news.Get("/api", controllers.APIGetNews)
	news.Get("/api/detail/:id", controllers.APIGetHiringDetailID)
	news.Post("/api/filter", controllers.APIPostNewsFilter)
	news.Post("/apply", controllers.PostNewsApply)
	news.Post("/cancel-apply", controllers.PostNewsCancelApply)

	receivedInvitation := app.Group("/received-inv", IsAuthenticated, CheckSession, CheckVerify, CheckPermissionGetInvitation)
	receivedInvitation.Get("", controllers.GetReceiveInvitationPage)
	receivedInvitation.Post("/api/filter", controllers.APIPostReceivedInvFilter)
	receivedInvitation.Post("/accept", controllers.PostReceivedInvAccept)
	//receivedInvitation.Post("/refuse", controllers.PostReceivedInvRefuse)

	errors := app.Group("/errors", IsAuthenticated, CheckSession)
	errors.Get("/403", controllers.GetError403)
	errors.Get("/401", controllers.GetError401)
	errors.Get("/404", controllers.GetError404)
}
