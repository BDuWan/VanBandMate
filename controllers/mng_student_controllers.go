package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"time"
)

func GetMngStudents(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngStudents")
	return c.Render("pages/managements/students/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func APIPostMngStudents(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngStudents")
	var users []models.User
	var query *gorm.DB
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	query = DB.Where("type_user_id", 5).Where(
		"deleted", false).Where(
		"referral_code", userLogin.CodeUser).Where(
		"verify", true).Where(
		"state", true)

	if userLogin.TypeUserID == 1 {
		query = DB.Where("type_user_id", 5).Where("deleted", false)
	}

	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

func GetMngStudentID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngStudentID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	var user models.User
	userId := c.Params("id")

	DB := initializers.DB

	if err := DB.Where("deleted", false).Where("type_user_id", 5).First(&user, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/managements/students/detail", fiber.Map{
		"User": user,
		"Ctx":  c,
	}, "layouts/main")
}

func APIPostAssignmentStudentID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAssignmentCourseID")

	var assignmentUser []models.AssignmentUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userId := c.Params("id")

	query = DB.Model(&models.AssignmentUser{}).Joins(
		"Assignment").Joins("Assignment.Course").Where(
		"Assignment.deleted", false).Where(
		"Assignment__Course.deleted", false).Where(
		"assignment_users.deleted", false).Where(
		"assignment_users.user_id", userId)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&assignmentUser).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "Assignment.name"
		}
	case 2:
		{
			sortColumn = "Assignment__Course.name"
		}
	case 3:
		{
			sortColumn = "Assignment.start_time"
		}
	case 4:
		{
			sortColumn = "Assignment.end_time"
		}
	default:
		{
			sortColumn = "Assignment.created_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("Assignment.name LIKE ? "+
			"or Assignment__Course.name LIKE ? "+
			"or Assignment.start_time LIKE ? "+
			"or Assignment.end_time LIKE ? ", search, search, search, search)
	}

	query.Find(&assignmentUser).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&assignmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            assignmentUser,
	})
}

func APIPostMngCoursesStudentID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngCoursesStudentID")
	var courses []models.CourseUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	userId := c.Params("id")

	query = DB.Model(&models.CourseUser{}).Joins("User").Joins("Course").Joins("Course.StudyProgram").Where("course_users.deleted", false).Where("course_users.user_id", userId)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&courses).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "Course.name"
		}
	case 2:
		{
			sortColumn = "Course__StudyProgram.title"
		}
	case 3:
		{
			sortColumn = "User.last_name"
		}
	case 4:
		{
			sortColumn = "Course.description"
		}
	default:
		{
			sortColumn = "Course.name"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("Course.name LIKE ? "+
			"or Course.description LIKE ? "+
			"or User.last_name LIKE ? "+
			"or Course__StudyProgram.title LIKE ? ", search, search, search, search)
	}

	query.Find(&courses).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            courses,
	})
}

func PostConfirmPaymentStudent(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostConfirmPaymentStudent")
	var user models.User
	DB := initializers.DB
	userId := c.Params("id")

	if err := DB.Where("user_id", userId).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found student ID: " + err.Error())
		return c.JSON("Can not find this student.")
	}

	user.Paid = true

	user.UpdatedAt = time.Now()

	if err := DB.Model(&user).Updates(user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course:  " + err.Error())
		return c.JSON("Cannot confirm payment for this student.")
	}
	var student_period models.StudentPeriod
	if err := DB.Where("student_id", userId).First(&student_period).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found student period: " + err.Error())
		return c.JSON("Success1")
	}
	if user.ReferralCode != "" {
		var resultGainCommissionSaleBusiness = GainCommissionSaleBusiness(user.ReferralCode, student_period.PeriodID)
		if resultGainCommissionSaleBusiness != "ok" {
			return c.JSON(resultGainCommissionSaleBusiness)
		}
	}

	return c.JSON("Success")
}
