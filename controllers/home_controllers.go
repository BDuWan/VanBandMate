package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"time"
)

//func GetHome(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetHome")
//	var studyPrograms []models.StudyProgram
//	var studyProgramUser []models.StudyProgramUser
//	DB := initializers.DB
//	userLogin := GetSessionUser(c)
//
//	if userLogin.RoleID== 4 {
//		var programId []int
//		if err := DB.Model(&models.CourseUser{}).Joins("Course").Where(
//			"Course.deleted", false).Where(
//			"course_users.deleted", false).Where(
//			"course_users.user_id", userLogin.UserID).Pluck("Course.study_program_id", &programId).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Where("deleted", false).Where("study_program_id in ?", programId).Find(&studyPrograms).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//	} else {
//		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//	}
//
//	if err := DB.Model(&models.StudyProgramUser{}).Joins("User").Joins("StudyProgram").Where(
//		"study_program_users.user_id", GetSessionUser(c).UserID).Where(
//		"StudyProgram.deleted", false).Where(
//		"User.deleted", false).Where(
//		"study_program_users.deleted", false).Find(&studyProgramUser).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	if userLogin.TypeUserID == 2 || userLogin.TypeUserID == 3 {
//		return c.Render("pages/home/index", fiber.Map{
//			"TypeUserID": userLogin.TypeUserID,
//			"Ctx":        c,
//		}, "layouts/main")
//	} else {
//		return c.Render("pages/home/index", fiber.Map{
//			"RoleID":       userLogin.RoleID,
//			"StudyProgram":     studyPrograms,
//			"StudyProgramUser": studyProgramUser,
//			"Ctx":              c,
//		}, "layouts/main")
//	}
//}

func GetHome(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetHome")
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	user := new(models.User)
	if err := DB.Model(&models.User{}).Joins("Role").Where("user_id", userLogin.UserID).First(user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Cannot get user")
	}
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	return c.Render("pages/home/index", fiber.Map{
		"Permissions": permissions,
		"User":        user,
		"Ctx":         c,
	}, "layouts/main")
}
