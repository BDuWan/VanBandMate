package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"io"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"os"
	"strconv"
	"time"
)

func GetMngStudyPrograms(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngStudyPrograms")
	return c.Render("pages/managements/study_programs/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func GetCreateMngStudyProgram(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateMngStudyProgram")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	return c.Render("pages/managements/study_programs/create", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func APIPostMngStudyPrograms(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngStudyPrograms")

	var studyPrograms []models.StudyProgram
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	query = DB.Where("deleted", false)

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
			sortColumn = "title"
		}
	case 3:
		{
			sortColumn = "description"
		}
	case 4:
		{
			sortColumn = "max_number"
		}
	case 5:
		{
			sortColumn = "number_student"
		}
	default:
		{
			sortColumn = "title"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("title LIKE ? "+
			"or number_student LIKE ? "+
			"or max_number LIKE ? "+
			"or description LIKE ? ", search, search, search, search)
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

func DeleteMngStudyProgramID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngStudyProgramID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.JSON("Permission Denied")
	}
	var StudyProgram models.StudyProgram
	var StudyProgramUser models.StudyProgramUser
	var course models.Course
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	StudyProgramId := c.Params("id")

	if err := DB.First(&StudyProgram, StudyProgramId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found StudyProgram ID: " + err.Error())
		return c.JSON("Can not delete this StudyProgram.")
	}

	StudyProgram.Deleted = true
	StudyProgram.DeletedBy = userLogin.UserID
	StudyProgram.DeletedAt = time.Now()

	if err := DB.Model(&StudyProgram).Updates(StudyProgram).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update StudyProgram:  " + err.Error())
		return c.JSON("Can not delete this StudyProgram.")
	}

	StudyProgramUser.Deleted = true
	StudyProgramUser.DeletedBy = userLogin.UserID
	StudyProgramUser.DeletedAt = time.Now()
	if err := DB.Model(&models.StudyProgramUser{}).Where(
		"study_program_id = ?", StudyProgramId).Updates(&StudyProgramUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not delete StudyProgram in StudyProgram user" + err.Error())
		return c.JSON("Can not delete this StudyProgram.")
	}

	course.Deleted = true
	course.DeletedBy = userLogin.UserID
	course.DeletedAt = time.Now()

	if err := DB.Model(&course).Where("study_program_id", StudyProgram.StudyProgramID).Updates(&course).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update StudyProgram:  " + err.Error())
		return c.JSON("Can not delete this course in StudyProgram.")
	}

	return c.JSON("Success")
}

func PostCreateMngStudyProgram(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateMngStudyProgram")

	var StudyProgram models.StudyProgram
	userLogin := GetSessionUser(c)
	if userLogin.TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	form, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	// Access the form fields
	title := c.FormValue("title")
	description := c.FormValue("description")
	numberMax, _ := strconv.Atoi(c.FormValue("number_max"))

	// Access the files
	files := form.File["image"]

	for _, file := range files {
		// Open the uploaded file
		src, err := file.Open()
		if err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
		defer src.Close()

		// Create a new file on the server
		dst, err2 := os.Create("public/assets/images/study_programs/" + file.Filename)
		if err2 != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
		defer dst.Close()

		// Copy the uploaded file content to the new file
		if _, err = io.Copy(dst, src); err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}

		StudyProgram.Title = title
		StudyProgram.Description = description
		StudyProgram.Image = file.Filename
		StudyProgram.MaxNumber = numberMax
		StudyProgram.Deleted = false
		StudyProgram.CreatedBy = userLogin.UserID
		StudyProgram.DeletedAt = time.Now()
		StudyProgram.CreatedAt = time.Now()
		if err := initializers.DB.Create(&StudyProgram).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.RedirectBack("")
		}

	}
	var program models.StudyProgram
	var courses []models.Course
	if err := initializers.DB.Where("deleted", false).Order("created_at asc").First(&program).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := initializers.DB.Where("deleted", false).Where("study_program_id", program.StudyProgramID).Find(&courses).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	for i := range courses {
		var course models.Course
		course.Name = courses[i].Name
		course.UserID = 0
		course.StudyProgramID = StudyProgram.StudyProgramID
		course.Description = ""
		course.Image = courses[i].Image
		course.Deleted = false
		course.CreatedBy = userLogin.UserID
		course.CreatedAt = time.Now()
		if err := initializers.DB.Create(&course).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}
	return c.Redirect("/managements/mng-study-programs")
}

func GetMngStudyProgramID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngStudyProgramID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	var StudyProgram models.StudyProgram

	DB := initializers.DB
	StudyProgramId := c.Params("id")

	err := DB.Where("deleted", false).First(&StudyProgram, StudyProgramId).Error
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	return c.Render("pages/managements/study_programs/edit", fiber.Map{
		"StudyProgram": StudyProgram,
		"Ctx":          c,
	}, "layouts/main")
}

func UpdateMngStudyProgramID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngStudyProgramID")
	var StudyProgram models.StudyProgram
	userLogin := GetSessionUser(c)
	if userLogin.TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	DB := initializers.DB
	StudyProgramId := c.Params("id")

	if err := DB.Where("study_program_id = ?", StudyProgramId).First(&StudyProgram).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	form, err := c.MultipartForm()
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	// Access the form fields
	title := c.FormValue("title")
	description := c.FormValue("description")
	numberMax, _ := strconv.Atoi(c.FormValue("number_max"))
	// Access the files
	files := form.File["image"]
	if len(files) == 0 {
		StudyProgram.Title = title
		StudyProgram.Description = description
		StudyProgram.MaxNumber = numberMax
		StudyProgram.UpdatedBy = userLogin.UserID
		StudyProgram.UpdatedAt = time.Now()

		if err := DB.Model(&StudyProgram).Updates(&StudyProgram).Error; err != nil {
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
			dst, err2 := os.Create("public/assets/images/study_programs/" + file.Filename)
			if err2 != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}
			defer dst.Close()

			// Copy the uploaded file content to the new file
			if _, err = io.Copy(dst, src); err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}

			StudyProgram.Title = title
			StudyProgram.Description = description
			StudyProgram.Image = file.Filename
			StudyProgram.MaxNumber = numberMax
			StudyProgram.UpdatedBy = userLogin.UserID
			StudyProgram.UpdatedAt = time.Now()

			if err := DB.Model(&StudyProgram).Updates(&StudyProgram).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
				return c.RedirectBack("")
			}
		}
	}

	return c.Redirect("/managements/mng-study-programs")
}
