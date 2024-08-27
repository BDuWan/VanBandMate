package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"time"
)

func GetFindMusicPlayerPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetFindMusicPlayerPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	if err := DB.Find(&districts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	if err := DB.Find(&wards).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	return c.Render("pages/find-music-player/index", fiber.Map{
		"This":           8,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"User":           userLogin,
		"provinces":      provinces,
		"districts":      districts,
		"wards":          wards,
		"Ctx":            c,
	}, "layouts/main")
}
