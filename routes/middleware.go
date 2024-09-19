package routes

import (
	"github.com/gofiber/fiber/v2"
	"vanbandmate/controllers"
)

func CheckAuthentication(c *fiber.Ctx) bool {
	sess, _ := controllers.SessAuth.Get(c)
	return sess.Get("login_success") != nil
}

func IsAuthenticated(c *fiber.Ctx) error {
	if !CheckAuthentication(c) {
		return c.Redirect("/login")
	}
	return c.Next()
}

func CheckSession(c *fiber.Ctx) error {
	userLogin := controllers.GetSessionUser(c)
	sess, _ := controllers.SessAuth.Get(c)

	if userLogin.Email != sess.Get("email") {
		return c.Redirect("/logout")
	}

	return c.Next()
}

func CheckVerify(c *fiber.Ctx) error {
	userLogin := controllers.GetSessionUser(c)

	if userLogin.Verify == false {
		return c.Redirect("/errors/401")
	}
	return c.Next()
}

func CheckPermissionMngUser(c *fiber.Ctx) error {
	if !controllers.CheckPermission("ql_tai_khoan", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionMngRole(c *fiber.Ctx) error {
	if !controllers.CheckPermission("ql_vai_tro", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionMngContract(c *fiber.Ctx) error {
	if !controllers.CheckPermission("ql_hop_dong", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionMngHiringNews(c *fiber.Ctx) error {
	if !controllers.CheckPermission("ql_tin_tuyen_dung", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionDashboard(c *fiber.Ctx) error {
	if !controllers.CheckPermission("thong_ke", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionMyContract(c *fiber.Ctx) error {
	if !controllers.CheckPermission("hop_dong_cua_toi", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionHiring(c *fiber.Ctx) error {
	if !controllers.CheckPermission("tuyen_dung", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionSendInvitation(c *fiber.Ctx) error {
	if !controllers.CheckPermission("loi_moi_da_gui", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionNews(c *fiber.Ctx) error {
	if !controllers.CheckPermission("xem_tin_tuyen_dung", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPermissionGetInvitation(c *fiber.Ctx) error {
	if !controllers.CheckPermission("loi_moi_nhan_duoc", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}
