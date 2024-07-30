package controllers

import (
	"lms/initializers"
	"lms/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
)

func CheckPermission(permissionName string, c *fiber.Ctx) bool {
	sess, _ := SessAuth.Get(c)
	userInterface := sess.Get("rolePermission")
	if userInterface == nil {
		outputdebug.String("[LMS]: " + "role permission not found in session")
		return false
	}

	rolePer, ok := userInterface.([]models.RolePermission)
	if !ok {
		outputdebug.String("[LMS]: " + "Failed to convert session data to role permission struct")
		return false
	}

	for _, per := range rolePer {
		if per.Permission.Permission == permissionName {
			return true
		}
	}

	return false

}

func GetSessionUser(c *fiber.Ctx) models.User {
	var user models.User
	var rolPer []models.RolePermission
	sess, _ := SessAuth.Get(c)
	email := sess.Get("email")

	if err := initializers.DB.Where("BINARY email = ?", email).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not found email in check session")
	}

	if err := initializers.DB.Model(&models.RolePermission{}).Preload("Permission").Where("role_id", user.RoleID).Find(&rolPer).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not found permission in check session")
	}
	//fmt.Println(rolPer)

	sess.Set("rolePermission", rolPer)

	if err := sess.Save(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not save session role permission in check session")
	}

	return user
}
