package routes

import (
	"lms/controllers"

	"github.com/gofiber/fiber/v2"
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

	if userLogin.Session != sess.Get("sessionId") {
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

func CheckPerHome(c *fiber.Ctx) error {
	if !controllers.CheckPermission("home", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerDashboard(c *fiber.Ctx) error {
	if !controllers.CheckPermission("dashboard", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerStudyPrograms(c *fiber.Ctx) error {
	if !controllers.CheckPermission("study_programs", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerCourses(c *fiber.Ctx) error {
	if !controllers.CheckPermission("courses", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngStudyPrograms(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_study_programs", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngCourses(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_courses", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngInstructors(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_instructors", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngStudents(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_students", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngSaleBusiness(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_sale_business", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerMngPayments(c *fiber.Ctx) error {
	if !controllers.CheckPermission("mng_payments", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerAccUsers(c *fiber.Ctx) error {
	if !controllers.CheckPermission("account_users", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}

func CheckPerAccRoles(c *fiber.Ctx) error {
	if !controllers.CheckPermission("account_roles", c) {
		return c.Redirect("/errors/403")
	}
	return c.Next()
}
