package controllers

import (
	"errors"
	"fmt"
	"io"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
)

func GetCourses(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCourses")
	var courses []models.Course
	var courseUsers []models.CourseUser

	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID == 5 || userLogin.TypeUserID == 4 {
		if err := DB.Model(&models.CourseUser{}).Joins("Course").Joins("Course.StudyProgram").Joins("User").Where(
			"course_users.deleted", false).Where(
			"Course.deleted", false).Where(
			"User.deleted", false).Where(
			"Course__StudyProgram.deleted", false).Where(
			"course_users.user_id", userLogin.UserID).Find(&courseUsers).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}

	} else {
		if err := DB.Model(&models.Course{}).Joins("StudyProgram").Where(
			"courses.deleted", false).Where(
			"StudyProgram.deleted", false).Where(
			"courses.user_id", userLogin.UserID).Find(&courses).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	return c.Render("pages/courses/index", fiber.Map{
		"courses":     courses,
		"CourseUsers": courseUsers,
		"TypeUserID":  userLogin.TypeUserID,
		"UserID":      userLogin.UserID,
		"Ctx":         c,
	}, "layouts/main")
}

//func GetCourseID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCourseID")
//	var course models.Course
//
//	courseId := c.Params("id")
//	DB := initializers.DB
//	userLogin := GetSessionUser(c)
//	var comments []models.Comment
//
//	if err := DB.Model(&models.Course{}).Joins("StudyProgram").Where(
//		"courses.deleted", false).Where(
//		"StudyProgram.deleted", false).Where(
//		"courses.course_id", courseId).First(&course).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	if err := DB.Model(&models.Comment{}).Joins("User").Where(
//		"comments.deleted", false).Where(
//		"User.deleted", false).Where(
//		"comments.course_id", courseId).Find(&comments).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	commentMap := make(map[int][]models.Comment)
//	parentComments := make(map[int]models.Comment)
//
//	for _, comment := range comments {
//		if comment.IsChildCmt {
//			commentMap[comment.ParentCmtID] = append(commentMap[comment.ParentCmtID], comment)
//		} else {
//			parentComments[comment.CommentID] = comment
//		}
//	}
//
//	var mapComments []models.Comment
//	for parentID, childComments := range commentMap {
//		if parentComment, exists := parentComments[parentID]; exists {
//			parentComment.SubComments = childComments
//			mapComments = append(mapComments, parentComment)
//		}
//	}
//
//	// Add parent comments without children to mapComments
//	for _, parentComment := range parentComments {
//		if _, exists := commentMap[parentComment.CommentID]; !exists {
//			mapComments = append(mapComments, parentComment)
//		}
//	}
//
//	return c.Render("pages/courses/detail", fiber.Map{
//		"UserLogin":         userLogin,
//		"UserLoginFullName": userLogin.GetFullName(),
//		"Course":            course,
//		"Comment":           mapComments,
//		"Ctx":               c,
//	}, "layouts/main")
//}

func GetCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCourseID")
	var course models.Course

	courseId := c.Params("id")
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if err := DB.Model(&models.Course{}).Joins("StudyProgram").Where(
		"courses.deleted", false).Where(
		"StudyProgram.deleted", false).Where(
		"courses.course_id", courseId).First(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.Render("pages/courses/detail", fiber.Map{
		"UserLogin": userLogin,
		//"UserLoginFullName": userLogin.GetFullName(),
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func APIPostStudentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostStudentCourseID")

	var courseUser []models.CourseUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	courseId := c.Params("id")

	query = DB.Model(&models.CourseUser{}).Joins(
		"User").Where(
		"course_users.deleted", false).Where(
		"User.deleted", false).Where(
		"User.type_user_id", 5).Where("course_users.course_id", courseId)

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
	query.Find(&courseUser).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "User.first_name"
		}
	case 2:
		{
			sortColumn = "User.last_name"
		}
	case 3:
		{
			sortColumn = "User.code_user"
		}
	case 4:
		{
			sortColumn = "User.address"
		}
	default:
		{
			sortColumn = "User.first_name"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("User.code_user LIKE ? "+
			"or User.address LIKE ? "+
			"or User.last_name LIKE ? "+
			"or User.first_name LIKE ? ", search, search, search, search)
	}

	query.Find(&courseUser).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&courseUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            courseUser,
	})
}

func APIPostDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostDocumentCourseID")

	var documents []models.Document
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	courseId := c.Params("id")

	query = DB.Model(&models.Document{}).Joins(
		"Course").Where(
		"Course.deleted", false).Where(
		"documents.deleted", false).Where(
		"documents.course_id", courseId)

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
	query.Find(&documents).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "documents.name"
		}
	case 2:
		{
			sortColumn = "documents.description"
		}
	default:
		{
			sortColumn = "documents.created_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("documents.name LIKE ? "+
			"or documents.description LIKE ? ", search, search)
	}

	query.Find(&documents).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&documents).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            documents,
	})
}

func DeleteDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteDocumentCourseID")
	var document models.Document

	userLogin := GetSessionUser(c)
	DB := initializers.DB
	documentId := c.Params("id")

	if err := DB.First(&document, documentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found document ID: " + err.Error())
		return c.JSON("Can not delete this document.")
	}

	document.Deleted = true
	document.DeletedBy = userLogin.UserID
	document.DeletedAt = time.Now()

	if err := DB.Model(&document).Updates(document).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update document:  " + err.Error())
		return c.JSON("Can not delete this document.")
	}

	return c.JSON("Success")
}

func GetCreateDocumentCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateDocumentCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/courses/documents/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func PostCreateDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateDocumentCourseID")
	var document models.Document
	userLogin := GetSessionUser(c)
	form, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Format failed")
	}

	// Access the form fields
	description := c.FormValue("description")
	courseId := c.Params("id")
	idCourse, err2 := strconv.Atoi(courseId)
	if err2 != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err2.Error())
		return c.JSON("Can not create document")
	}

	if _, err := os.Stat("documents"); os.IsNotExist(err) {
		os.Mkdir("documents", 0755)
	}

	if _, err := os.Stat("documents/" + courseId); os.IsNotExist(err) {
		os.Mkdir("documents/"+courseId, 0755)
	}

	// Access the files
	files := form.File["name"]

	for _, file := range files {
		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer src.Close()

		if file.Size > 1024*1024*100 {
			return c.JSON("File needs to be less than 100MB")

		}
		// Create a new file on the server

		dst, err2 := os.Create("documents/" + courseId + "/" + file.Filename)
		if err2 != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		if _, err = io.Copy(dst, src); err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}

		document.Name = file.Filename
		document.Description = description
		document.CourseID = idCourse
		document.CreatedBy = userLogin.UserID
		document.DeletedAt = time.Now()
		document.CreatedAt = time.Now()
		if err := initializers.DB.Create(&document).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}

	}
	return c.JSON("Success")

}

func GetDownloadDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetDownloadDocumentCourseID")

	var document models.Document

	if err := initializers.DB.First(&document, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	filePath := "./documents/" + strconv.Itoa(document.CourseID) + "/" + document.Name // Adjust the file path accordingly

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return c.RedirectBack("")
	}

	return c.Download(filePath, document.Name)

}

func GetEditDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditDocumentCourseID")

	var document models.Document

	if err := initializers.DB.First(&document, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/courses/documents/edit", fiber.Map{
		"Document": document,
		"Ctx":      c,
	}, "layouts/main")

}

func UpdateDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateDocumentCourseID")
	var document models.Document
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	form, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Format failed")
	}

	// Access the form fields
	description := c.FormValue("description")
	documentId := c.Params("id")

	if err := DB.First(&document, documentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found document")
	}

	if _, err := os.Stat("documents"); os.IsNotExist(err) {
		os.Mkdir("documents", 0755)
	}

	courseId := strconv.Itoa(document.CourseID)
	if _, err := os.Stat("documents/" + courseId); os.IsNotExist(err) {
		os.Mkdir("documents/"+courseId, 0755)
	}

	// Access the files
	files := form.File["name"]
	if len(files) == 0 {
		document.Description = description
		document.UpdatedAt = time.Now()
		document.UpdatedBy = userLogin.UserID
		if err := DB.Model(&document).Updates(&document).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update document")
			return c.JSON("Can not update document")
		}
		return c.JSON("Success")

	}

	for _, file := range files {
		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer src.Close()

		if file.Size > 1024*1024*100 {
			return c.JSON("File needs to be less than 100MB")

		}
		// Create a new file on the server

		dst, err2 := os.Create("documents/" + courseId + "/" + file.Filename)
		if err2 != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		if _, err = io.Copy(dst, src); err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}

		document.Name = file.Filename
		document.Description = description
		document.UpdatedBy = userLogin.UserID
		document.UpdatedAt = time.Now()
		if err := DB.Model(&document).Updates(&document).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update document")
			return c.JSON("Can not update document")
		}

	}
	return c.JSON("Success")

}

func APIPostLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostLessonCourseID")

	var lessons []models.Lesson
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	courseId := c.Params("id")

	query = DB.Model(&models.Lesson{}).Joins(
		"Course").Where(
		"Course.deleted", false).Where(
		"lessons.deleted", false).Where(
		"lessons.course_id", courseId)

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
	query.Find(&lessons).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "lessons.name"
		}
	case 2:
		{
			sortColumn = "lessons.start_time"
		}
	case 3:
		{
			sortColumn = "lessons.end_time"
		}
	case 4:
		{
			sortColumn = "lessons.description"
		}
	case 5:
		{
			sortColumn = "lessons.link_study"
		}
	case 6:
		{
			sortColumn = "lessons.link_record"
		}
	default:
		{
			sortColumn = "lessons.created_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("lessons.name LIKE ? "+
			"or lessons.start_time LIKE ? "+
			"or lessons.description LIKE ? "+
			"or lessons.end_time LIKE ? "+
			"or lessons.link_study LIKE ? "+
			"or lessons.link_record LIKE ? ", search, search, search, search, search, search)
	}

	query.Find(&lessons).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&lessons).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            lessons,
	})
}

func DeleteLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteLessonCourseID")
	var lesson models.Lesson

	userLogin := GetSessionUser(c)
	DB := initializers.DB
	lessonId := c.Params("id")

	if err := DB.First(&lesson, lessonId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found lesson ID: " + err.Error())
		return c.JSON("Can not delete this lesson.")
	}

	lesson.Deleted = true
	lesson.DeletedBy = userLogin.UserID
	lesson.DeletedAt = time.Now()

	if err := DB.Model(&lesson).Updates(lesson).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update lesson:  " + err.Error())
		return c.JSON("Can not delete this lesson.")
	}

	return c.JSON("Success")
}

func GetCreateLessonCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateLessonCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/courses/lessons/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func PostCreateLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateLessonCourseID")
	var lesson models.Lesson

	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.Lesson)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.StartTime == "" {
		return c.JSON("Haven't chosen a start time")
	}
	if form.EndTime == "" {
		return c.JSON("Haven't chosen a end time")
	}

	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
	if errStr != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
	}
	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
	if errEnd != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
	}

	if startTime.After(endTime) {
		return c.JSON("The end time must be after the start time")
	}

	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
	if errOff != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
	}
	lesson.Name = form.Name
	lesson.StartTime = startTime.Add(offset)
	lesson.EndTime = endTime.Add(offset)
	lesson.Description = form.Description
	lesson.LinkStudy = form.LinkStudy
	lesson.LinkRecord = form.LinkRecord
	lesson.CourseID = form.CourseID
	lesson.Deleted = false
	lesson.CreatedBy = userLogin.UserID
	lesson.DeletedAt = time.Now()
	lesson.CreatedAt = time.Now()

	if err := DB.Create(&lesson).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not create lesson")
	}

	var userCourse []string

	if err := DB.Table("course_users a").Where("a.deleted = false").
		Joins("JOIN users b on a.user_id = b.user_id AND (b.deleted = false and b.type_user_id = 5)").
		Where("a.course_id", lesson.CourseID).Group("b.email").
		Pluck("b.email", &userCourse).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	for _, item := range userCourse {
		body := fmt.Sprintf(os.Getenv("HOST")+"/courses/lesson/%s", strconv.Itoa(lesson.LessonID))
		go SendEmail("New lesson", "Your have new lesson: "+lesson.Name+".\n"+"Join the lesson at: "+body, item)
	}
	return c.JSON("Success")

}

func GetEditLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditLessonCourseID")

	var lesson models.Lesson

	if err := initializers.DB.First(&lesson, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/courses/lessons/edit", fiber.Map{
		"Lesson": lesson,
		"Ctx":    c,
	}, "layouts/main")

}
func GetInfoLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetInfoLessonCourseID")

	var lesson models.Lesson

	if err := initializers.DB.First(&lesson, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/courses/lessons/info", fiber.Map{
		"Lesson": lesson,
		"Ctx":    c,
	}, "layouts/main")

}

func UpdateLessonCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateLessonCourseID")
	var lesson models.Lesson

	lessonId := c.Params("id")
	DB := initializers.DB

	form := new(structs.Lesson)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if err := DB.Where("deleted", false).First(&lesson, lessonId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found lesson")
	}
	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.StartTime == "" {
		return c.JSON("Haven't chosen a start time")
	}
	if form.EndTime == "" {
		return c.JSON("Haven't chosen a end time")
	}
	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
	if errStr != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
	}
	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
	if errEnd != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
	}
	if startTime.After(endTime) {
		return c.JSON("The end time must be after the start time")
	}

	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
	if errOff != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
	}
	lesson.Name = form.Name
	lesson.StartTime = startTime.Add(offset)
	lesson.EndTime = endTime.Add(offset)
	lesson.Description = form.Description
	lesson.LinkStudy = form.LinkStudy
	lesson.LinkRecord = form.LinkRecord
	lesson.CourseID = form.CourseID
	lesson.UpdatedAt = time.Now()
	lesson.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Model(&lesson).Updates(&lesson).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update lesson")
		return c.JSON("Can not update lesson")
	}

	return c.JSON("Success")

}

func APIPostAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAssignmentCourseID")

	var assignments []models.Assignment
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	courseId := c.Params("id")

	query = DB.Model(&models.Assignment{}).Joins(
		"Course").Where(
		"Course.deleted", false).Where(
		"assignments.deleted", false).Where(
		"assignments.course_id", courseId)

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
	query.Find(&assignments).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "assignments.name"
		}
	case 2:
		{
			sortColumn = "assignments.description"
		}
	case 3:
		{
			sortColumn = "assignments.start_time"
		}
	case 4:
		{
			sortColumn = "assignments.end_time"
		}
	default:
		{
			sortColumn = "assignments.created_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("assignments.name LIKE ? "+
			"or assignments.start_time LIKE ? "+
			"or assignments.description LIKE ? "+
			"or assignments.end_time LIKE ? ", search, search, search, search)
	}

	query.Find(&assignments).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&assignments).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            assignments,
	})
}

func DeleteAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteAssignmentCourseID")
	var assignment models.Assignment
	var assignmentUser models.AssignmentUser

	userLogin := GetSessionUser(c)
	DB := initializers.DB
	assignmentId := c.Params("id")

	if err := DB.First(&assignment, assignmentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found lesson ID: " + err.Error())
		return c.JSON("Can not delete this assignment.")
	}

	assignment.Deleted = true
	assignment.DeletedBy = userLogin.UserID
	assignment.DeletedAt = time.Now()

	if err := DB.Model(&assignment).Updates(assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update assignment:  " + err.Error())
		return c.JSON("Can not delete this assignment.")
	}

	assignmentUser.Deleted = true
	assignmentUser.DeletedBy = userLogin.UserID
	assignmentUser.DeletedAt = time.Now()

	if err := DB.Model(&models.AssignmentUser{}).Where(
		"assignment_id = ?", assignmentId).Updates(&assignmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not delete assignment in user assignment" + err.Error())
		return c.JSON("Can not delete this assignment.")
	}

	return c.JSON("Success")
}

func GetCreateAssignmentCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateAssignmentCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/courses/assignments/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func PostCreateAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateAssignmentCourseID")
	var assignment models.Assignment
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.Assignment)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}
	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.StartTime == "" {
		return c.JSON("Haven't chosen a start time")
	}
	if form.EndTime == "" {
		return c.JSON("Haven't chosen a end time")
	}

	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
	if errStr != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
	}
	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
	if errEnd != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
	}

	if startTime.After(endTime) {
		return c.JSON("The end time must be after the start time")
	}

	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
	if errOff != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
	}
	assignment.Name = form.Name
	assignment.StartTime = startTime.Add(offset)
	assignment.EndTime = endTime.Add(offset)
	assignment.CourseID = form.CourseID
	assignment.Description = form.Description
	assignment.Deleted = false
	assignment.CreatedBy = userLogin.UserID
	assignment.DeletedAt = time.Now()
	assignment.CreatedAt = time.Now()

	if err := DB.Create(&assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not create assignment")
	}

	var courseUser []models.CourseUser
	if err := DB.Where("deleted", false).Where("course_id", assignment.CourseID).Find(&courseUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not create assignment")
	}

	for i := range courseUser {
		var assignUser models.AssignmentUser
		assignUser.AssignmentID = assignment.AssignmentID
		assignUser.UserID = courseUser[i].UserID
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

	var userAssign []string

	if err := DB.Table("assignment_users a").Where("a.deleted = false").
		Joins("JOIN users b on a.user_id = b.user_id AND (b.deleted = false and b.type_user_id = 5)").
		Joins("JOIN assignments c on a.assignment_id = c.assignment_id AND c.deleted = false").
		Where("c.course_id", assignment.CourseID).Group("b.email").
		Pluck("b.email", &userAssign).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	for _, item := range userAssign {
		body := fmt.Sprintf(os.Getenv("HOST")+"/courses/assignment/%s", strconv.Itoa(assignment.AssignmentID))
		go SendEmail("New assignment", "Your have new assignment: "+assignment.Name+".\n"+"Join the assignment at: "+body, item)
	}
	return c.JSON("Success")

}

func GetEditAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditAssignmentCourseID")

	var assignment models.Assignment

	if err := initializers.DB.First(&assignment, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/courses/assignments/edit", fiber.Map{
		"Assignment": assignment,
		"Ctx":        c,
	}, "layouts/main")

}

func UpdateAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateAssignmentCourseID")
	var assignment models.Assignment

	assignmentId := c.Params("id")
	DB := initializers.DB

	form := new(structs.Assignment)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if err := DB.Where("deleted", false).First(&assignment, assignmentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found assignment")
	}
	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.CourseID == 0 {
		return c.JSON("Haven't chosen a course")
	}
	if form.StartTime == "" {
		return c.JSON("Haven't chosen a start time")
	}
	if form.EndTime == "" {
		return c.JSON("Haven't chosen a end time")
	}

	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
	if errStr != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
	}
	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
	if errEnd != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
	}

	if startTime.After(endTime) {
		return c.JSON("The end time must be after the start time")
	}

	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
	if errOff != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
	}
	assignment.Name = form.Name
	assignment.StartTime = startTime.Add(offset)
	assignment.EndTime = endTime.Add(offset)
	assignment.CourseID = form.CourseID
	assignment.Description = form.Description
	assignment.UpdatedAt = time.Now()
	assignment.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Model(&assignment).Updates(&assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update assignment")
		return c.JSON("Can not update assignment")
	}

	return c.JSON("Success")

}

func GetAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAssignmentCourseID")
	var assignment models.Assignment
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	if err := DB.Where(
		"assignments.deleted", false).Where(
		"assignments.assignment_id", c.Params("id")).First(&assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if userLogin.TypeUserID == 5 {

		var assignmentUser models.AssignmentUser

		status := 1
		result := ""
		//var assignmentCourse []structs.AssignmentCourse

		if err := initializers.DB.Where(
			"user_id", userLogin.UserID).Where(
			"assignment_id", c.Params("id")).First(&assignmentUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		}

		if assignmentUser.Status > 0 {
			status = assignmentUser.Status
		}
		if assignmentUser.Result != "" {
			result = assignmentUser.Result
		}

		return c.Render("pages/courses/assignments/detail", fiber.Map{
			"Status":     status,
			"Result":     result,
			"Assignment": assignment,
			"TypeUserID": userLogin.TypeUserID,
			"Ctx":        c,
		}, "layouts/main")
	}

	return c.Render("pages/courses/assignments/list", fiber.Map{
		"Assignment": assignment,
		"TypeUserID": userLogin.TypeUserID,
		"Ctx":        c,
	}, "layouts/main")
}

func GetAssignmentUserCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAssignmentUserCourseID")
	var assignment models.Assignment
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	if err := DB.Where(
		"assignments.deleted", false).Where(
		"assignments.assignment_id", c.Params("id")).First(&assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	var assignmentUser models.AssignmentUser

	status := 1
	result := ""
	//var assignmentCourse []structs.AssignmentCourse

	if err := initializers.DB.Where(
		"user_id", c.Params("uid")).Where(
		"assignment_id", c.Params("id")).First(&assignmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

	}

	if assignmentUser.Status > 0 {
		status = assignmentUser.Status
	}
	if assignmentUser.Result != "" {
		result = assignmentUser.Result
	}

	return c.Render("pages/courses/assignments/instr_detail", fiber.Map{
		"Status":     status,
		"UserID":     assignmentUser.UserID,
		"Result":     result,
		"Assignment": assignment,
		"TypeUserID": userLogin.TypeUserID,
		"Ctx":        c,
	}, "layouts/main")

}

func APIPostFileUserAssignmentID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostFileUserAssignmentID")

	var files []models.FileUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	query = DB.Where(
		"deleted", false).Where(
		"assignment_id", c.Params("id")).Where(
		"user_id", GetSessionUser(c).UserID)

	if userLogin.TypeUserID != 5 {
		query = DB.Where(
			"deleted", false).Where(
			"assignment_id", c.Params("id")).Where(
			"user_id", c.Params("uid"))
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
	query.Find(&files).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "file"
		}
	default:
		{
			sortColumn = "created_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("file LIKE ? ", search)
	}

	query.Find(&files).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&files).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            files,
	})
}

func PostFileUserAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostFileUserAssignmentCourseID")
	var assignmentUser models.AssignmentUser
	userLogin := GetSessionUser(c)
	form, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Format failed")
	}

	// Access the form fields
	assignmentId := c.FormValue("assignment_id")

	idAssign, err3 := strconv.Atoi(assignmentId)
	if err3 != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err3.Error())
	}

	if err := initializers.DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Where(
		"assignment_id", assignmentId).First(&assignmentUser).Error; err != nil {

	}

	if assignmentUser.Status != 1 {
		return c.JSON("Can not upload file. Please rollback assignment.")
	}

	if _, err := os.Stat("documents"); os.IsNotExist(err) {
		os.Mkdir("documents", 0755)
	}
	if _, err := os.Stat("documents/assignments"); os.IsNotExist(err) {
		os.Mkdir("documents/assignments", 0755)
	}

	if _, err := os.Stat("documents/assignments/" + assignmentId); os.IsNotExist(err) {
		os.Mkdir("documents/assignments/"+assignmentId, 0755)
	}

	if _, err := os.Stat("documents/assignments/" + assignmentId + "/" + strconv.Itoa(userLogin.UserID)); os.IsNotExist(err) {
		os.Mkdir("documents/assignments/"+assignmentId+"/"+strconv.Itoa(userLogin.UserID), 0755)
	}

	// Access the files
	files := form.File["files[]"]

	for _, file := range files {
		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer src.Close()

		if file.Size > 1024*1024*100 {
			return c.JSON("File needs to be less than 100MB")

		}
		// Create a new file on the server

		dst, err2 := os.Create("documents/assignments/" + assignmentId + "/" + strconv.Itoa(userLogin.UserID) + "/" + file.Filename)
		if err2 != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		if _, err = io.Copy(dst, src); err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create document")
		}
		var fileUser models.FileUser
		fileUser.AssignmentID = idAssign
		fileUser.UserID = userLogin.UserID
		fileUser.File = file.Filename
		fileUser.Deleted = false
		fileUser.CreatedBy = userLogin.UserID
		fileUser.DeletedAt = time.Now()
		fileUser.CreatedAt = time.Now()
		if err := initializers.DB.Create(&fileUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create fileUser")
		}

	}

	return c.JSON("Success")

}

func PostSubmitAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostSubmitAssignmentCourseID")
	var assigmentUser models.AssignmentUser
	var assignment models.Assignment
	var status = 1
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	assignmentId := c.FormValue("assignment_id")
	currentTime := time.Now()
	if err := DB.First(&assignment, assignmentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found file ID: " + err.Error())
		return c.JSON("Can not delete this file.")
	}

	if assignment.EndTime.After(currentTime) {
		status = 2
	} else if assignment.EndTime.Before(currentTime) {
		status = 3
	}

	idAssign, err3 := strconv.Atoi(assignmentId)
	if err3 != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err3.Error())
	}

	if err := DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Where(
		"assignment_id", assignmentId).First(&assigmentUser).Error; err != nil {
		var assign models.AssignmentUser
		assign.AssignmentID = idAssign
		assign.UserID = userLogin.UserID
		assign.Result = ""
		assign.Status = status
		assign.Deleted = false
		assign.CreatedBy = userLogin.UserID
		assign.DeletedAt = time.Now()
		assign.CreatedAt = time.Now()
		if err := DB.Create(&assign).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not create assigmentUser")
		}
		return c.JSON(status)
	}
	assigmentUser.Status = status
	assigmentUser.UpdatedBy = userLogin.UserID
	assigmentUser.UpdatedAt = time.Now()
	if err := DB.Model(&assigmentUser).Updates(&assigmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not update assigmentUser")
	}

	return c.JSON(status)

}

func DeleteFileUserCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteFileUserCourseID")
	var file models.FileUser
	var assignmentUser models.AssignmentUser
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	fileId := c.Params("id")

	if err := DB.First(&file, fileId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found file ID: " + err.Error())
		return c.JSON("Can not delete this file.")
	}

	if err := initializers.DB.Where("deleted", false).Where(
		"user_id", userLogin.UserID).Where(
		"assignment_id", file.AssignmentID).First(&assignmentUser).Error; err != nil {

	}

	if assignmentUser.Status != 1 {
		return c.JSON("Can not delete this file. Please rollback assignment.")
	}

	file.Deleted = true
	file.DeletedBy = userLogin.UserID
	file.DeletedAt = time.Now()

	if err := DB.Model(&file).Updates(file).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update file:  " + err.Error())
		return c.JSON("Can not delete this file.")
	}

	return c.JSON("Success")
}

func GetDownloadFileUserCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetDownloadFileUserCourseID")

	var file models.FileUser

	if err := initializers.DB.First(&file, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	filePath := "./documents/assignments/" + strconv.Itoa(file.AssignmentID) + "/" + strconv.Itoa(file.UserID) + "/" + file.File // Adjust the file path accordingly

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return c.RedirectBack("")
	}

	return c.Download(filePath, file.File)

}

func UpdateRollbackAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateRollbackAssignmentCourseID")
	var assigmentUser models.AssignmentUser
	var user models.User
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	assignmentId := c.FormValue("assignment_id")
	if c.Params("id") == "1" {
		if err := DB.Where("deleted", false).Where(
			"user_id", userLogin.UserID).Where(
			"assignment_id", assignmentId).First(&assigmentUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not found assigmentUser")
		}
	} else {
		if err := DB.Where("deleted", false).Where(
			"user_id", c.Params("id")).Where(
			"assignment_id", assignmentId).First(&assigmentUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not found assigmentUser")
		}
	}

	assigmentUser.Status = 1
	assigmentUser.UpdatedBy = userLogin.UserID
	assigmentUser.UpdatedAt = time.Now()
	if err := DB.Model(&assigmentUser).Updates(&assigmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not update assigmentUser")
	}

	if err := DB.Select("email").Where("deleted", false).Where("type_user_id = 5").First(
		&user, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found user")
	}
	body := "The instructor has rollback your assignment. Access to assignment " + os.Getenv("HOST") + "/courses/assignment/" + assignmentId
	go SendEmail("Result assignment", body, user.Email)

	return c.JSON("Success")

}

func UpdateResultAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateResultAssignmentCourseID")
	var assigmentUser models.AssignmentUser
	var user models.User
	DB := initializers.DB
	userId := c.FormValue("user_id")
	assignmentId := c.FormValue("assignment_id")
	result := c.FormValue("result")

	if err := DB.Where("deleted", false).Where(
		"user_id", userId).Where(
		"assignment_id", assignmentId).First(&assigmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found assigmentUser")
	}
	assigmentUser.Result = result
	assigmentUser.UpdatedBy = GetSessionUser(c).UserID
	assigmentUser.UpdatedAt = time.Now()
	if err := DB.Model(&assigmentUser).Updates(&assigmentUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not update assigmentUser")
	}

	if err := DB.Select("email").Where("deleted", false).Where("type_user_id = 5").First(
		&user, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found user")
	}
	body := "The instructor has commented on your assignment. Access to assignment " + os.Getenv("HOST") + "/courses/assignment/" + assignmentId
	go SendEmail("Result assignment", body, user.Email)

	return c.JSON("Success")

}

func APIPostListUserAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostListUserAssignmentCourseID")

	var assignmentUser []models.AssignmentUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	query = DB.Model(&models.AssignmentUser{}).Joins("User").Where(
		"assignment_users.deleted", false).Where(
		"User.deleted", false).Where(
		"assignment_users.assignment_id", c.Params("id"))

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
			sortColumn = "User.first_name"
		}
	case 2:
		{
			sortColumn = "User.last_name"
		}
	default:
		{
			sortColumn = "assignment_users.updated_at"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("User.last_name LIKE ? or User.first_name LIKE ?", search, search)
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

func GetCreateCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateCourse")
	var studyPrograms []models.StudyProgramUser
	DB := initializers.DB

	if err := DB.Model(&models.StudyProgramUser{}).Joins("StudyProgram").Where(
		"study_program_users.deleted", false).Where(
		"study_program_users.user_id", c.Params("id")).Find(&studyPrograms).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
	}
	return c.Render("pages/courses/create", fiber.Map{
		"StudyPrograms": studyPrograms,
		"UserID":        c.Params("id"),
		"Ctx":           c,
	}, "layouts/main")
}

func PostCreateCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateCourse")
	var course models.Course

	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.Course)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.StudyProgramID == 0 {
		return c.JSON("Haven't chosen a study program")
	}
	if form.UserID == 0 {
		return c.JSON("Haven't chosen a instructor")
	}
	if form.Description == "" {
		return c.JSON("Description has not been entered ")
	}

	course.Name = form.Name
	//course.StudyProgramID = form.StudyProgramID
	course.StudyProgramID = form.StudyProgramID
	course.UserID = form.UserID
	course.Description = form.Description
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

func GetEditCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditCourse")
	var course models.Course
	var studyPrograms []models.StudyProgramUser

	DB := initializers.DB
	courseId := c.Params("id")
	if err := DB.Where("deleted", false).First(&course, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Model(&models.StudyProgramUser{}).Joins("StudyProgram").Where(
		"study_program_users.deleted", false).Where(
		"study_program_users.user_id", course.UserID).Find(&studyPrograms).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
	}

	return c.Render("pages/courses/edit", fiber.Map{
		"StudyPrograms": studyPrograms,
		"Course":        course,
		"Ctx":           c,
	}, "layouts/main")
}

func UpdateCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateCourse")
	var course models.Course

	courseId := c.Params("id")
	DB := initializers.DB

	form := new(structs.Course)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if err := DB.Where("deleted", false).First(&course, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not found course")
	}
	if form.Name == "" {
		return c.JSON("Name has not been entered")
	}
	if form.StudyProgramID == 0 {
		return c.JSON("Haven't chosen a study program")
	}
	if form.UserID == 0 {
		return c.JSON("Haven't chosen a instructor")
	}
	if form.Description == "" {
		return c.JSON("Description has not been entered ")
	}

	course.Name = form.Name
	//course.StudyProgramID = form.StudyProgramID
	course.StudyProgramID = form.StudyProgramID
	course.UserID = form.UserID
	course.Description = form.Description
	course.UpdatedAt = time.Now()
	course.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Model(&course).Updates(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update course")
		return c.JSON("Can not update course")
	}

	return c.JSON("Success")
}
