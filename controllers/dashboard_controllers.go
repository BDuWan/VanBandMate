package controllers

//
//import (
//	"github.com/gofiber/fiber/v2"
//	"github.com/zetamatta/go-outputdebug"
//	"lms/initializers"
//	"lms/models"
//	"lms/structs"
//	"time"
//)
//
//func GetDashboard(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetDashboard")
//
//	var students, instructors, saleBusiness []models.User
//	var topCourseUser []structs.TopCourseUser
//	var topSaleUser []structs.TopSaleUser
//	var courseUsers []models.CourseUser
//	var money int
//
//	DB := initializers.DB
//
//	userLogin := GetSessionUser(c)
//
//	if userLogin.TypeUserID == 1 {
//		var pricePro models.PriceProgram
//		if err := DB.Select("user_id").Where("type_user_id", 5).Where("deleted", false).Find(&students).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Select("user_id").Where("type_user_id", 4).Where("deleted", false).Find(&instructors).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Where("type_user_id > 1").Where("type_user_id < 4").Where("deleted", false).Find(&saleBusiness).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Where("deleted", false).Order("price_program_id desc").First(&pricePro).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		money = len(students) * pricePro.Price
//
//		if err := DB.Select("Course.course_id as course_id", "Course.title as title", "count(course_users.course_id) as number").Table("course_users").Model(&models.CourseUser{}).Joins(
//			"Course").Joins(
//			"User").Where(
//			"Course.deleted", false).Where(
//			"User.deleted", false).Where(
//			"course_users.deleted", false).Where(
//			"User.type_user_id = 5").Group(
//			"course_users.course_id").Order(
//			"number desc").Limit(5).Find(&topCourseUser).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//		if err := DB.Select("referral_code", "count(user_id) as number").Table("users").Where(
//			"deleted", false).Where(
//			"referral_code <> ''").Where(
//			"type_user_id", 5).Group(
//			"referral_code").Order(
//			"number desc").Limit(5).Find(&topSaleUser).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//		if err := DB.Model(&models.CourseUser{}).Joins("User").Joins("Course").Where(
//			"Course.deleted", false).Where(
//			"User.deleted", false).Where(
//			"course_users.deleted", false).Where(
//			"type_user_id", 5).Order(
//			"course_users.created_at desc").Limit(8).Find(&courseUsers).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//	} else {
//		var priceU models.PriceUser
//		if err := DB.Select("user_id").Where("type_user_id", 5).Where("deleted", false).Where("referral_code", userLogin.CodeUser).Find(&students).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Select("user_id").Where("type_user_id", 4).Where("deleted", false).Where("referral_code", userLogin.CodeUser).Find(&instructors).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Where("type_user_id > 1").Where("type_user_id < 4").Where("code_user", userLogin.CodeUser).Where("deleted", false).Find(&saleBusiness).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		if err := DB.Model(&models.PriceUser{}).Joins("PriceProgram").Where("price_users.deleted", false).Where("user_id", userLogin.UserID).First(&priceU).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//		money = len(students) * priceU.PriceProgram.Commission
//
//		if err := DB.Select("Course.course_id as course_id", "Course.title as title", "count(course_users.course_id) as number").Table("course_users").Model(&models.CourseUser{}).Joins(
//			"Course").Joins(
//			"User").Where(
//			"Course.deleted", false).Where(
//			"User.deleted", false).Where(
//			"course_users.deleted", false).Where(
//			"User.type_user_id = 5").Where("User.referral_code", userLogin.CodeUser).Group(
//			"course_users.course_id").Order(
//			"number desc").Limit(5).Find(&topCourseUser).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//		if err := DB.Select("referral_code", "count(user_id) as number").Table("users").Where(
//			"deleted", false).Where(
//			"type_user_id", 5).Where(
//			"referral_code", userLogin.CodeUser).Order(
//			"number desc").Limit(5).Find(&topSaleUser).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//		if err := DB.Model(&models.CourseUser{}).Joins("User").Joins("Course").Where(
//			"Course.deleted", false).Where(
//			"User.deleted", false).Where(
//			"course_users.deleted", false).Where(
//			"User.type_user_id", 5).Where("User.referral_code", userLogin.CodeUser).Order(
//			"course_users.created_at desc").Limit(8).Find(&courseUsers).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//
//	}
//	for i := range topSaleUser {
//		for j := range saleBusiness {
//			if saleBusiness[j].CodeUser == topSaleUser[i].ReferralCode {
//				if saleBusiness[j].TypeUserID == 2 {
//					topSaleUser[i].Name = saleBusiness[j].FirstName + " " + saleBusiness[j].LastName
//				} else {
//					topSaleUser[i].Name = saleBusiness[j].NameBusiness
//				}
//
//			}
//		}
//	}
//	return c.Render("pages/dashboard/index", fiber.Map{
//		"StudentC":      len(students),
//		"InstructorC":   len(instructors),
//		"SaleBusinessC": len(saleBusiness),
//		"Money":         money,
//		"TopCourseUser": topCourseUser,
//		"TopSaleUser":   topSaleUser,
//		"CourseUsers":   courseUsers,
//		"Ctx":           c,
//	}, "layouts/main")
//}
//
//func Notification(c *fiber.Ctx) error {
//
//	var assignment []models.AssignmentUser
//	var lessons []models.Lesson
//
//	userLogin := GetSessionUser(c)
//	DB := initializers.DB
//
//	if err := DB.Model(&models.AssignmentUser{}).Preload("Assignment").Preload("Assignment.Course").Where("status", 1).Where("deleted", false).Where("user_id", userLogin.UserID).Find(&assignment).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	if err := DB.Where("deleted", false).Where("course_id IN (select course_id from course_users where user_id = ? and deleted = ?)", userLogin.UserID, false).Find(&lessons).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	var count = 0
//
//	var lessonArr []int
//	currentTime := time.Now()
//	for i := range lessons {
//		startTime := lessons[i].StartTime
//		endTime := lessons[i].EndTime
//
//		if startTime.Before(currentTime) && endTime.After(currentTime) {
//			count = count + 1
//			lessonArr = append(lessonArr, lessons[i].LessonID)
//		}
//
//	}
//
//	if err := DB.Model(&models.Lesson{}).Preload("Course").Where("deleted", false).Where("lesson_id IN ?", lessonArr).Find(&lessons).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"AssignmentN": len(assignment),
//		"LessonN":     count,
//		"Assignments": assignment,
//		"Lessons":     lessons,
//	})
//
//}
