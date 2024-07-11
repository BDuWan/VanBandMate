package controllers

import (
	"fmt"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
)

func GetMngInstructors(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngInstructors")
	return c.Render("pages/managements/instructors/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func APIPostMngInstructors(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngInstructors")
	var users []models.User
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	query = DB.Where("type_user_id", 4).Where(
		"deleted", false).Where(
		"referral_code", userLogin.CodeUser).Where(
		"verify", true).Where(
		"state", true)

	if userLogin.TypeUserID == 1 {
		query = DB.Where("type_user_id", 4).Where(
			"deleted", false).Where(
			"verify", true).Where(
			"state", true)
	}

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
	query.Find(&users).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "first_name"
		}
	case 2:
		{
			sortColumn = "last_name"
		}
	case 3:
		{
			sortColumn = "email"
		}
	case 4:
		{
			sortColumn = "phone_number"
		}
	case 5:
		{
			sortColumn = "address"
		}
	default:
		{
			sortColumn = "first_name"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("first_name LIKE ? "+
			"or last_name LIKE ? "+
			"or phone_number LIKE ? "+
			"or email LIKE ? "+
			"or address LIKE ? ", search, search, search, search, search)
	}

	query.Find(&users).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            users,
	})
}

func GetMngInstructorID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngInstructorID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	var user models.User
	userId := c.Params("id")

	DB := initializers.DB

	if err := DB.Where("deleted", false).Where("type_user_id", 4).First(&user, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/managements/instructors/detail", fiber.Map{
		"User": user,
		"Ctx":  c,
	}, "layouts/main")
}

func APIPostMngStudyProgramsInstructorID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngStudyProgramsInstructorID")
	var studyPrograms []models.StudyProgramUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	userId := c.Params("id")

	query = DB.Model(&models.StudyProgramUser{}).Joins("StudyProgram").Where(
		"study_program_users.user_id", userId).Where(
		"StudyProgram.deleted", false).Where(
		"study_program_users.deleted", false)

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
	query.Find(&studyPrograms).Count(&totalRecords)

	switch req.Order[0].Column {
	case 2:
		{
			sortColumn = "StudyProgram.title"
		}
	case 4:
		{
			sortColumn = "StudyProgram.description"
		}
	default:
		{
			sortColumn = "StudyProgram.title"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("StudyProgram.title LIKE ? "+
			"or StudyProgram.description LIKE ? ", search, search)
	}

	query.Find(&studyPrograms).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&studyPrograms).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            studyPrograms,
	})
}

func DeleteMngInstructorStudyProgramID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngInstructorStudyProgramID")
	var studyProgramUser models.StudyProgramUser
	//var course models.Course
	var courseUsers []models.CourseUser
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	studyProgramId := c.Params("id")

	if err := DB.Where("study_program_user_id", studyProgramId).First(&studyProgramUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_program ID: " + err.Error())
		return c.JSON("Can not delete this study_program in detail instructor.")
	}

	studyProgramUser.Deleted = true
	studyProgramUser.DeletedBy = userLogin.UserID
	studyProgramUser.DeletedAt = time.Now()

	if err := DB.Model(&studyProgramUser).Updates(&studyProgramUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update study_program:  " + err.Error())
		return c.JSON("Can not delete this study_program.")
	}

	if err := DB.Model(&models.CourseUser{}).Joins("Course").Joins("Course.StudyProgram").Where(
		"Course.deleted", false).Where(
		"Course__StudyProgram.deleted", false).Where(
		"course_users.deleted", false).Where(
		"course_users.user_id", studyProgramUser.UserID).Where(
		"Course.study_program_id", studyProgramUser.StudyProgramID).Find(&courseUsers).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not find course user of this study program:  " + err.Error())
		return c.JSON("Can not find course user of this study program.")
	}
	for _, courseUser := range courseUsers {
		courseUser.Deleted = true
		courseUser.DeletedBy = userLogin.UserID
		courseUser.DeletedAt = time.Now()
		if err := DB.Updates(&courseUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course user of this study program:  " + err.Error())
			return c.JSON("Can not update course user of this study program.")
		}
	}

	return c.JSON("Success")
}

func APIPostMngCoursesInstructorID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngCoursesInstructorID")
	var Courses []models.CourseUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	userId := c.Params("id")

	query = DB.Model(&models.CourseUser{}).Joins("Course").Joins("Course.StudyProgram").Where(
		"Course.deleted", false).Where(
		"Course__StudyProgram.deleted", false).Where(
		"course_users.deleted", false).Where("course_users.user_id", userId)

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
	query.Find(&Courses).Count(&totalRecords)

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
			"or Course__StudyProgram.title LIKE ? ", search, search, search)
	}

	query.Find(&Courses).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&Courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            Courses,
	})
}

func DeleteMngInstructorCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngInstructorCourseID")
	var course models.CourseUser
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	courseUserId := c.Params("id")

	if err := DB.Where("course_user_id", courseUserId).First(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course ID: " + err.Error())
		return c.JSON("Can not delete this course in detail instructor.")
	}

	course.Deleted = true
	course.DeletedBy = userLogin.UserID
	course.DeletedAt = time.Now()

	if err := DB.Model(&course).Updates(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course:  " + err.Error())
		return c.JSON("Can not delete this course.")
	}

	return c.JSON("Success")
}

func GetCreateMngInstructorsStudyProgram(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateMngInstructorsStudyProgram")
	var studyPrograms []models.StudyProgram
	DB := initializers.DB
	userId := c.Params("id")
	if err := DB.Where("study_program_id NOT IN (SELECT study_program_id FROM study_program_users WHERE deleted = false and user_id = ?)", userId).Where("deleted", false).Find(&studyPrograms).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
	}

	return c.Render("pages/managements/instructors/create_study_program", fiber.Map{
		"UserID":        userId,
		"StudyPrograms": studyPrograms,
		"Ctx":           c,
	}, "layouts/main")
}

func PostCreateMngInstructorsStudyProgram(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateMngInstructorsStudyProgram")
	var studyProgram models.StudyProgramUser
	var studyProgramCheck models.StudyProgramUser
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.StudyProgramInstructor)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}
	if form.UserID == 0 {
		return c.JSON("Haven't a user")
	}
	if form.StudyProgramID == 0 {
		return c.JSON("Haven't chosen a study_program")
	}
	if err := DB.Where("study_program_id", form.StudyProgramID).Where(
		"user_id", form.UserID).Where(
		"deleted", false).First(&studyProgramCheck).Error; err != nil {
		studyProgram.StudyProgramID = form.StudyProgramID
		studyProgram.UserID = form.UserID
		studyProgram.Deleted = false
		studyProgram.CreatedBy = userLogin.UserID
		studyProgram.DeletedAt = time.Now()
		studyProgram.CreatedAt = time.Now()

		if err := DB.Create(&studyProgram).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create study program")
		}

		return c.JSON("Success")
	}
	return c.JSON("Duplicate study program cannot be added")
}

func GetCreateMngInstructorsCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateMngInstructorsCourse")
	var courses []models.Course
	DB := initializers.DB

	var program []int

	if err := DB.Table("study_program_users").Where(
		"deleted", false).Where("user_id", c.Params("id")).Pluck("study_program_id", &program).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in create course:  " + err.Error())
	}

	if err := DB.Model(&models.Course{}).Joins("StudyProgram").Where(
		"StudyProgram.deleted", false).Where(
		"courses.deleted", false).Where(
		"StudyProgram.study_program_id in ?", program).Find(&courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in create course:  " + err.Error())
	}

	return c.Render("pages/managements/instructors/create_course", fiber.Map{
		"Courses": courses,
		"UserID":  c.Params("id"),
		"Ctx":     c,
	}, "layouts/main")
}

func PostCreateMngInstructorsCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateMngInstructorsCourse")
	var course, courseCheck models.CourseUser

	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.CourseInstructor)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.UserID == 0 {
		return c.JSON("Haven't chosen a user")
	}

	if err := DB.Where("course_id", form.CourseID).Where(
		"user_id", form.UserID).Where(
		"deleted", false).First(&courseCheck).Error; err != nil {
		course.CourseID = form.CourseID
		course.UserID = form.UserID
		course.Deleted = false
		course.CreatedBy = userLogin.UserID
		course.DeletedAt = time.Now()
		course.CreatedAt = time.Now()

		if err := DB.Create(&course).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create course")
		}

		return c.JSON("Success")

	}
	return c.JSON("Duplicate study program cannot be added")
}

func GetEditMngInstructorCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditMngInstructorCourse")
	var courses []models.Course
	var course models.CourseUser
	DB := initializers.DB
	var program []int
	if err := DB.Where(
		"deleted", false).First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in create course:  " + err.Error())
	}
	if err := DB.Table("study_program_users").Where(
		"deleted", false).Where("user_id", course.UserID).Pluck("study_program_id", &program).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in create course:  " + err.Error())
	}

	if err := DB.Model(&models.Course{}).Joins("StudyProgram").Where(
		"StudyProgram.deleted", false).Where(
		"courses.deleted", false).Where(
		"StudyProgram.study_program_id in ?", program).Find(&courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in create course:  " + err.Error())
	}

	return c.Render("pages/managements/instructors/edit_course", fiber.Map{
		"Courses": courses,
		"Course":  course,
		"Ctx":     c,
	}, "layouts/main")
}

func UpdateMngInstructorCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngInstructorCourse")
	var course models.CourseUser

	courseUserId := c.Params("id")
	DB := initializers.DB

	form := new(structs.CourseUser)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if err := DB.Where("deleted", false).First(&course, courseUserId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found course")
	}
	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.UserID == 0 {
		return c.JSON("Haven't chosen a user")
	}

	course.CourseID = form.CourseID
	course.UpdatedAt = time.Now()
	course.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Model(&course).Updates(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update course")
		return c.JSON("Can not update course")
	}

	return c.JSON("Success")
}

func GetSelectMngInstructorsCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetSelectMngInstructorsCourse")
	var courses []models.Course
	DB := initializers.DB

	if err := DB.Where("deleted", false).Where(
		"study_program_id", c.Params("id")).Find(&courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course in program:  " + err.Error())
	}

	return c.JSON(courses)
}
