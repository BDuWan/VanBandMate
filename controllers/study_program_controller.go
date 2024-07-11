package controllers

import (
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
)

func GetStudyPrograms(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetStudyPrograms")
	var studyPrograms []models.StudyProgramUser
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID < 4 {
		return c.Redirect("/errors/404")
	}

	if err := DB.Model(&models.StudyProgramUser{}).Joins("User").Joins("StudyProgram").Where(
		"study_program_users.user_id", userLogin.UserID).Where(
		"StudyProgram.deleted", false).Where(
		"User.deleted", false).Where(
		"study_program_users.deleted", false).Find(&studyPrograms).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	return c.Render("pages/study_programs/index", fiber.Map{
		"StudyPrograms": studyPrograms,
		"Ctx":           c,
	}, "layouts/main")
}

func GetStudyProgramID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetStudyProgramID")
	var studyProgram models.StudyProgram
	var courses []models.Course
	var courseUsers []models.CourseUser
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID < 4 {
		return c.Redirect("/errors/404")
	}
	if !userLogin.Paid {
		return c.Redirect("/errors/404")
	}

	if err := DB.Where("deleted", false).Where("study_program_id", c.Params("id")).First(&studyProgram).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	if userLogin.TypeUserID == 5 {
		if err := DB.Model(&models.Course{}).Where(
			"courses.deleted", false).Where(
			"courses.study_program_id", c.Params("id")).Find(&courses).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	} else {
		if err := DB.Model(&models.Course{}).Where(
			"courses.deleted", false).Where(
			"courses.study_program_id", c.Params("id")).Where(
			"courses.user_id", userLogin.UserID).Find(&courses).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	if err := DB.Model(&models.CourseUser{}).Joins("Course").Joins("User").Where(
		"course_users.deleted", false).Where(
		"Course.deleted", false).Where(
		"User.deleted", false).Where(
		"course_users.user_id", userLogin.UserID).Find(&courseUsers).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.Render("pages/study_programs/detail", fiber.Map{
		"TypeUserID":   userLogin.TypeUserID,
		"StudyProgram": studyProgram,
		"Courses":      courses,
		"CourseUsers":  courseUsers,
		"Ctx":          c,
	}, "layouts/main")
}

func PostCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCourse")
	var courseUser, check models.CourseUser

	userLogin := GetSessionUser(c)
	if userLogin.TypeUserID != 5 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found user login")
		return c.JSON("Not found user login")
	}
	DB := initializers.DB

	form := new(structs.CourseUser)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}
	// fmt.Println(form.CourseID)
	// fmt.Println(userLogin.UserID)
	if err := DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Where(
		"course_id", form.CourseID).First(&check).Error; err != nil {

		courseUser.CourseID = form.CourseID
		courseUser.UserID = userLogin.UserID
		courseUser.Deleted = false
		courseUser.CreatedBy = userLogin.UserID
		courseUser.DeletedAt = time.Now()
		courseUser.CreatedAt = time.Now()

		if err := DB.Create(&courseUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not add course")
		}

		// gain := GainNumberStudent(form.CourseID)
		// if gain != "ok" {
		// 	return c.JSON(gain)
		// }

		var assignments []models.Assignment
		if err := DB.Where("deleted", false).Where("course_id", courseUser.CourseID).Find(&assignments).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create assignment")
		}

		for i := range assignments {
			var assignUser models.AssignmentUser
			assignUser.AssignmentID = assignments[i].AssignmentID
			assignUser.UserID = courseUser.UserID
			assignUser.Status = 1
			assignUser.Deleted = false
			assignUser.CreatedBy = userLogin.UserID
			assignUser.DeletedAt = time.Now()
			assignUser.CreatedAt = time.Now()
			if err := DB.Create(&assignUser).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
				return c.JSON("Can not create assignment user")
			}
		}

		return c.JSON("Success")
	}

	return c.JSON("This course has been taken")
}

func GainNumberStudent(studyProgramId int) string {
	DB := initializers.DB
	var StudyProgram models.StudyProgram
	if err := DB.Where("deleted", false).Where(
		"study_program_id", studyProgramId).First(&StudyProgram).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Can not find study program"
	}
	if StudyProgram.NumberStudent >= StudyProgram.MaxNumber {
		return "Out of slots"
	}

	StudyProgram.NumberStudent += 1
	if err := DB.Model(&StudyProgram).Updates(&StudyProgram).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Can not gain study program number student"
	}
	return "ok"
}
