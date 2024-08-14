package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"time"
)

func GetMyContractPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get Profile ID")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	return c.Render("pages/info/profile", fiber.Map{
		"Permissions": permissions,
		"User":        userLogin,
		"Ctx":         c,
	}, "layouts/main")
}
