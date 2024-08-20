package controllers

//import (
//	"github.com/gofiber/fiber/v2"
//	"github.com/zetamatta/go-outputdebug"
//	"lms/initializers"
//	"lms/models"
//	"time"
//)
//
//func GetHome(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetHome")
//	DB := initializers.DB
//	userLogin := GetSessionUser(c)
//	user := new(models.User)
//	if err := DB.Model(&models.User{}).Joins("Role").Where("user_id", userLogin.UserID).First(user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Cannot get user")
//	}
//	sess, _ := SessAuth.Get(c)
//	permissions := sess.Get("rolePermission")
//
//	return c.Render("pages/info/index", fiber.Map{
//		"Permissions": permissions,
//		"User":        user,
//		"Ctx":         c,
//	}, "layouts/main")
//}
