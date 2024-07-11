package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"time"
)

func GetHome(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetHome")
	var studyPrograms []models.StudyProgram
	var studyProgramUser []models.StudyProgramUser
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID == 4 {
		var programId []int
		if err := DB.Model(&models.CourseUser{}).Joins("Course").Where(
			"Course.deleted", false).Where(
			"course_users.deleted", false).Where(
			"course_users.user_id", userLogin.UserID).Pluck("Course.study_program_id", &programId).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
		if err := DB.Where("deleted", false).Where("study_program_id in ?", programId).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}

	} else {
		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	if err := DB.Model(&models.StudyProgramUser{}).Joins("User").Joins("StudyProgram").Where(
		"study_program_users.user_id", GetSessionUser(c).UserID).Where(
		"StudyProgram.deleted", false).Where(
		"User.deleted", false).Where(
		"study_program_users.deleted", false).Find(&studyProgramUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	if userLogin.TypeUserID == 2 || userLogin.TypeUserID == 3 {
		return c.Render("pages/home/index_sale", fiber.Map{
			"TypeUserID": userLogin.TypeUserID,
			"Ctx":        c,
		}, "layouts/main")
	} else {
		return c.Render("pages/home/index", fiber.Map{
			"TypeUserID":       userLogin.TypeUserID,
			"StudyProgram":     studyPrograms,
			"StudyProgramUser": studyProgramUser,
			"Ctx":              c,
		}, "layouts/main")
	}

}

func PostStudyProgram(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostStudyProgram")
	var studyProgramUser, check models.StudyProgramUser
	var checkExist []models.StudyProgramUser

	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID != 5 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found user login")
		return c.JSON("Not found user login")
	}

	DB := initializers.DB

	form := new(structs.StudyProgramUser)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if err := DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Find(&checkExist).Error; err != nil {
		return c.JSON("An error occurred while retrieving data")
	}

	if len(checkExist) > 0 {
		return c.JSON("You can only participate in 1 study program")
	}

	if err := DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Where(
		"study_program_id", form.StudyProgramID).First(&check).Error; err != nil {

		studyProgramUser.StudyProgramID = form.StudyProgramID
		studyProgramUser.UserID = userLogin.UserID
		studyProgramUser.Deleted = false
		studyProgramUser.CreatedBy = userLogin.UserID
		studyProgramUser.DeletedAt = time.Now()
		studyProgramUser.CreatedAt = time.Now()

		gain := GainNumberStudent(form.StudyProgramID)
		if gain != "ok" {
			return c.JSON(gain)
		}

		if err := DB.Create(&studyProgramUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not add study program")
		}
		return c.JSON("Success")
	}
	return c.JSON("This study program has been taken")
}
