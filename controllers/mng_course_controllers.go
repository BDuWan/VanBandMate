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

func GetMngCourses(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCourses")
	return c.Render("pages/managements/courses/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func GetCreateMngCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateMngCourse")
	var studyPrograms []models.StudyProgram
	var users []models.User
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID == 1 {
		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
		}

		if err := DB.Where("deleted", false).Where("type_user_id", 4).Find(&users).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found users in create course:  " + err.Error())
		}
	} else {
		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
		}

		if err := DB.Where("deleted", false).Where("referral_code", userLogin.CodeUser).Where("type_user_id", 4).Find(&users).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found users in create course:  " + err.Error())
		}
	}

	return c.Render("pages/managements/courses/create", fiber.Map{
		"StudyPrograms": studyPrograms,
		"Users":         users,
		"Ctx":           c,
	}, "layouts/main")
}

func APIPostMngCourses(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngCourses")
	var courses []models.Course
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	userLogin := GetSessionUser(c)
	if userLogin.TypeUserID == 1 {
		query = DB.Model(&models.Course{}).Joins("StudyProgram").Where("courses.deleted", false)
	} else {
		query = DB.Model(&models.Course{}).Joins("StudyProgram").Where("courses.deleted", false).Where("courses.user_id", userLogin.CodeUser)
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
	query.Find(&courses).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "courses.image"
		}
	case 2:
		{
			sortColumn = "courses.name"
		}
	case 3:
		{
			sortColumn = "StudyProgram.title"
		}
	case 4:
		{
			sortColumn = "courses.price"
		}
	case 5:
		{
			sortColumn = "courses.description"
		}
	default:
		{
			sortColumn = "courses.name"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("courses.name LIKE ? "+
			"or courses.description LIKE ? "+
			"or StudyProgram.title LIKE ? "+
			"or courses.price LIKE ? ", search, search, search, search)
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

func DeleteMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngCourseID")
	var course models.Course
	var courseUser models.CourseUser
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	courseId := c.Params("id")

	if err := DB.First(&course, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course ID: " + err.Error())
		return c.JSON("Can not delete this course.")
	}

	course.Deleted = true
	course.DeletedBy = userLogin.UserID
	course.DeletedAt = time.Now()

	if err := DB.Model(&course).Updates(course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course:  " + err.Error())
		return c.JSON("Can not delete this course.")
	}

	courseUser.Deleted = true
	courseUser.DeletedBy = userLogin.UserID
	courseUser.DeletedAt = time.Now()
	if err := DB.Model(&models.CourseUser{}).Where(
		"course_id = ?", courseId).Updates(&courseUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not delete course in course user" + err.Error())
		return c.JSON("Can not delete this course.")
	}

	return c.JSON("Success")
}

func PostCreateMngCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateMngCourse")
	var course models.Course

	userLogin := GetSessionUser(c)

	DB := initializers.DB

	formFile, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	files := formFile.File["image"]
	name := c.FormValue("name")
	studyProgramId, _ := strconv.Atoi(c.FormValue("study_program_id"))
	description := c.FormValue("description")

	for _, file := range files {
		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
		defer src.Close()

		// Create a new file on the server
		dst, err2 := os.Create("public/assets/images/courses/" + file.Filename)
		if err2 != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		if _, err = io.Copy(dst, src); err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}

		course.Name = name
		course.StudyProgramID = studyProgramId
		course.Image = file.Filename
		course.Description = description
		course.Deleted = false
		course.CreatedBy = userLogin.UserID
		course.DeletedAt = time.Now()
		course.CreatedAt = time.Now()

		if err := DB.Create(&course).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.RedirectBack("")
		}

	}

	return c.Redirect("/managements/mng-courses")

}

func GetMngEditCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngEditCourseID")
	var course models.Course
	var studyPrograms []models.StudyProgram
	var users []models.User

	DB := initializers.DB
	courseId := c.Params("id")
	if err := DB.Where("deleted", false).First(&course, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	userLogin := GetSessionUser(c)

	if userLogin.TypeUserID == 1 {
		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
		}

		if err := DB.Where("deleted", false).Where("type_user_id", 4).Find(&users).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found users in create course:  " + err.Error())
		}
	} else {
		if err := DB.Where("deleted", false).Find(&studyPrograms).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found study_programs in create course:  " + err.Error())
		}

		if err := DB.Where("deleted", false).Where("referral_code", userLogin.CodeUser).Where("type_user_id", 4).Find(&users).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found users in create course:  " + err.Error())
		}
	}

	return c.Render("pages/managements/courses/edit", fiber.Map{
		"Users":         users,
		"StudyPrograms": studyPrograms,
		"Course":        course,
		"Ctx":           c,
	}, "layouts/main")
}

func UpdateMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngCourseID")
	var course models.Course
	userLogin := GetSessionUser(c)
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

	formFile, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	files := formFile.File["image"]
	name := c.FormValue("name")
	studyProgramId, _ := strconv.Atoi(c.FormValue("study_program_id"))
	description := c.FormValue("description")

	if len(files) == 0 {
		course.Name = name
		course.StudyProgramID = studyProgramId
		course.Description = description
		course.UpdatedAt = time.Now()
		course.UpdatedBy = userLogin.UserID

		if err := DB.Model(&course).Updates(&course).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.RedirectBack("")
		}
	} else {
		for _, file := range files {
			// Open the uploaded file
			src, err := file.Open()
			if err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}
			defer src.Close()

			// Create a new file on the server
			dst, err2 := os.Create("public/assets/images/courses/" + file.Filename)
			if err2 != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}
			defer dst.Close()

			// Copy the uploaded file content to the new file
			if _, err = io.Copy(dst, src); err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}

			course.Name = name
			course.StudyProgramID = studyProgramId
			course.Image = file.Filename
			course.Description = description
			course.UpdatedAt = time.Now()
			course.UpdatedBy = userLogin.UserID

			if err := DB.Model(&course).Updates(&course).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
				return c.RedirectBack("")
			}
		}
	}

	return c.Redirect("/managements/mng-courses")

}

func GetMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCourseID")
	var course models.Course
	var courseUser models.CourseUser
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	courseId := c.Params("id")

	if err := DB.Model(&models.Course{}).Joins(
		"StudyProgram").Where(
		"courses.deleted", false).Where(
		"StudyProgram.deleted", false).First(&course, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := DB.Model(&models.CourseUser{}).Joins("User").Where(
		"course_users.deleted", false).First(&courseUser, courseId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.Render("pages/managements/courses/detail", fiber.Map{
		"UserLogin":  userLogin,
		"CourseUser": courseUser,
		"Course":     course,
		"Ctx":        c,
	}, "layouts/main")
}

func APIPostStudentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostStudentMngCourseID")

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
			sortColumn = "User.email"
		}
	case 4:
		{
			sortColumn = "User.phone_number"
		}
	case 5:
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
		query = query.Where("User.email LIKE ? "+
			"or User.phone_number LIKE ? "+
			"or User.address LIKE ? "+
			"or User.last_name LIKE ? "+
			"or User.first_name LIKE ? ", search, search, search, search, search)
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

func APIPostDocumentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostDocumentMngCourseID")

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
	case 3:
		{
			sortColumn = "documents.created_at"
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
			"or documents.description LIKE ? "+
			"or documents.created_at LIKE ? ", search, search, search)
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

func DeleteDocumentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteDocumentMngCourseID")
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

func GetCreateDocumentMngCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateDocumentMngCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/managements/courses/documents/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func PostCreateDocumentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateDocumentMngCourse")
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

func GetDownloadMngDocumentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetDownloadMngDocumentCourseID")

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

func GetEditDocumentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditDocumentMngCourseID")

	var document models.Document

	if err := initializers.DB.First(&document, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/managements/courses/documents/edit", fiber.Map{
		"Document": document,
		"Ctx":      c,
	}, "layouts/main")

}

func UpdateDocumentMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateDocumentMngCourseID")
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

func APIPostLessonMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostLessonMngCourseID")

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

func DeleteLessonMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteLessonMngCourseID")
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

func GetCreateLessonMngCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateLessonMngCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/managements/courses/lessons/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}

func PostCreateLessonMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateLessonMngCourseID")
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
	lesson.LinkRecord = form.LinkRecord
	lesson.LinkStudy = form.LinkStudy
	lesson.Description = form.Description
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

func GetEditLessonMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetEditLessonMngCourseID")

	var lesson models.Lesson

	if err := initializers.DB.First(&lesson, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/managements/courses/lessons/edit", fiber.Map{
		"Lesson": lesson,
		"Ctx":    c,
	}, "layouts/main")

}

func UpdateLessonMngCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateLessonMngCourseID")
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
	lesson.LinkRecord = form.LinkRecord
	lesson.LinkStudy = form.LinkStudy
	lesson.CourseID = form.CourseID
	lesson.UpdatedAt = time.Now()
	lesson.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Model(&lesson).Updates(&lesson).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update lesson")
		return c.JSON("Can not update lesson")
	}

	return c.JSON("Success")

}

func GetMngAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngAssignmentCourseID")
	var assignment models.Assignment
	DB := initializers.DB
	userLogin := GetSessionUser(c)
	if err := DB.Where(
		"assignments.deleted", false).Where(
		"assignments.assignment_id", c.Params("id")).First(&assignment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.Render("pages/managements/courses/assignments/list", fiber.Map{
		"Assignment": assignment,
		"TypeUserID": userLogin.TypeUserID,
		"Ctx":        c,
	}, "layouts/main")
}

func GetMngAssignmentUserCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngAssignmentUserCourseID")
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

	return c.Render("pages/managements/courses/assignments/instr_detail", fiber.Map{
		"Status":     status,
		"UserID":     assignmentUser.UserID,
		"Result":     result,
		"Assignment": assignment,
		"TypeUserID": userLogin.TypeUserID,
		"Ctx":        c,
	}, "layouts/main")

}

func GetMngEditAssignmentCourseID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngEditAssignmentCourseID")

	var assignment models.Assignment

	if err := initializers.DB.First(&assignment, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}

	return c.Render("pages/managements/courses/assignments/edit", fiber.Map{
		"Assignment": assignment,
		"Ctx":        c,
	}, "layouts/main")

}

func GetMngCreateAssignmentCourse(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCreateAssignmentCourse")
	var course models.Course

	if err := initializers.DB.First(&course, c.Params("id")).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.RedirectBack("")
	}
	return c.Render("pages/managements/courses/assignments/create", fiber.Map{
		"Course": course,
		"Ctx":    c,
	}, "layouts/main")
}
