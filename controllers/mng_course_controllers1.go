package controllers

// import (
// 	"fmt"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/zetamatta/go-outputdebug"
// 	"gorm.io/gorm"
// 	"io"
// 	"lms/initializers"
// 	"lms/models"
// 	"lms/structs"
// 	"os"
// 	"time"
// )

// func GetMngCourses(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCourses")
// 	return c.Render("pages/managements/courses/index", fiber.Map{
// 		"Ctx": c,
// 	}, "layouts/main")
// }

// func GetCreateMngCourse(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetCreateMngCourse")
// 	if GetSessionUser(c).TypeUserID > 1 {
// 		return c.Redirect("/errors/403")
// 	}
// 	return c.Render("pages/managements/courses/create", fiber.Map{
// 		"Ctx": c,
// 	}, "layouts/main")
// }

// func APIPostMngCourses(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngCourses")

// 	var courses []models.Course
// 	var query *gorm.DB
// 	var req structs.ReqBody
// 	DB := initializers.DB

// 	query = DB.Where("deleted", false)

// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Cannot parse JSON",
// 		})
// 	}

// 	var totalRecords int64
// 	var filteredRecords int64

// 	var sortColumn string
// 	var sortDir string

// 	sortDir = req.Order[0].Dir
// 	query.Find(&courses).Count(&totalRecords)

// 	switch req.Order[0].Column {
// 	case 2:
// 		{
// 			sortColumn = "courses.title"
// 		}
// 	case 3:
// 		{
// 			sortColumn = "courses.price"
// 		}
// 	case 4:
// 		{
// 			sortColumn = "courses.description"
// 		}
// 	default:
// 		{
// 			sortColumn = "courses.title"
// 			sortDir = "asc"
// 		}
// 	}
// 	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
// 	query = query.Order(orderBy)

// 	if req.Search.Value != "" {
// 		search := "%" + req.Search.Value + "%"
// 		query = query.Where("courses.title LIKE ? "+
// 			"or courses.price LIKE ? "+
// 			"or courses.description LIKE ? ", search, search, search)
// 	}

// 	query.Find(&courses).Count(&filteredRecords)

// 	query = query.Offset(req.Start).Limit(req.Length)

// 	if err := query.Find(&courses).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 	}

// 	return c.JSON(fiber.Map{
// 		"draw":            req.Draw,
// 		"recordsTotal":    totalRecords,
// 		"recordsFiltered": filteredRecords,
// 		"data":            courses,
// 	})
// }

// func DeleteMngCourseID(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngCourseID")
// 	if GetSessionUser(c).TypeUserID > 1 {
// 		return c.JSON("Permission Denied")
// 	}
// 	var course models.Course
// 	var courseUser models.CourseUser
// 	var class models.Class
// 	userLogin := GetSessionUser(c)
// 	DB := initializers.DB
// 	courseId := c.Params("id")

// 	if err := DB.First(&course, courseId).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found course ID: " + err.Error())
// 		return c.JSON("Can not delete this course.")
// 	}

// 	course.Deleted = true
// 	course.DeletedBy = userLogin.UserID
// 	course.DeletedAt = time.Now()

// 	if err := DB.Model(&course).Updates(course).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course:  " + err.Error())
// 		return c.JSON("Can not delete this course.")
// 	}

// 	courseUser.Deleted = true
// 	courseUser.DeletedBy = userLogin.UserID
// 	courseUser.DeletedAt = time.Now()
// 	if err := DB.Model(&models.CourseUser{}).Where(
// 		"course_id = ?", courseId).Updates(&courseUser).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not delete course in course user" + err.Error())
// 		return c.JSON("Can not delete this course.")
// 	}

// 	class.Deleted = true
// 	class.DeletedBy = userLogin.UserID
// 	class.DeletedAt = time.Now()

// 	if err := DB.Model(&class).Where("course_id", course.CourseID).Updates(&class).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update course:  " + err.Error())
// 		return c.JSON("Can not delete this class in course.")
// 	}

// 	return c.JSON("Success")
// }

// func PostCreateMngCourse(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostCreateMngCourse")

// 	var course models.Course
// 	userLogin := GetSessionUser(c)
// 	if userLogin.TypeUserID > 1 {
// 		return c.Redirect("/errors/403")
// 	}
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 	}

// 	// Access the form fields
// 	title := c.FormValue("title")
// 	description := c.FormValue("description")
// 	price := c.FormValue("price")

// 	// Access the files
// 	files := form.File["image"]

// 	for _, file := range files {
// 		// Open the uploaded file
// 		src, err := file.Open()
// 		if err != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		}
// 		defer src.Close()

// 		// Create a new file on the server
// 		dst, err2 := os.Create("public/assets/images/courses/" + file.Filename)
// 		if err2 != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		}
// 		defer dst.Close()

// 		// Copy the uploaded file content to the new file
// 		if _, err = io.Copy(dst, src); err != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		}

// 		course.Title = title
// 		course.Description = description
// 		course.Price = price
// 		course.Image = file.Filename
// 		course.Deleted = false
// 		course.CreatedBy = userLogin.UserID
// 		course.DeletedAt = time.Now()
// 		course.CreatedAt = time.Now()
// 		if err := initializers.DB.Create(&course).Error; err != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			return c.RedirectBack("")
// 		}

// 	}
// 	return c.Redirect("/managements/mng-courses")

// }

// func GetMngCourseID(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCourseID")
// 	if GetSessionUser(c).TypeUserID > 1 {
// 		return c.Redirect("/errors/403")
// 	}
// 	var course models.Course

// 	DB := initializers.DB
// 	courseId := c.Params("id")

// 	err := DB.Where("deleted", false).First(&course, courseId).Error
// 	if err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

// 		return c.RedirectBack("")
// 	}

// 	return c.Render("pages/managements/courses/edit", fiber.Map{
// 		"Course": course,
// 		"Ctx":    c,
// 	}, "layouts/main")
// }

// func UpdateMngCourseID(c *fiber.Ctx) error {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngCourseID")
// 	var course models.Course
// 	userLogin := GetSessionUser(c)
// 	if userLogin.TypeUserID > 1 {
// 		return c.Redirect("/errors/403")
// 	}
// 	DB := initializers.DB
// 	courseId := c.Params("id")

// 	if err := DB.Where("course_id = ?", courseId).First(&course).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

// 		return c.RedirectBack("")
// 	}

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 	}

// 	// Access the form fields
// 	title := c.FormValue("title")
// 	description := c.FormValue("description")
// 	price := c.FormValue("price")

// 	// Access the files
// 	files := form.File["image"]
// 	if len(files) == 0 {
// 		course.Title = title
// 		course.Description = description
// 		course.Price = price
// 		course.UpdatedBy = userLogin.UserID
// 		course.UpdatedAt = time.Now()

// 		if err := DB.Model(&course).Updates(&course).Error; err != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			return c.RedirectBack("")
// 		}
// 	} else {
// 		for _, file := range files {
// 			// Open the uploaded file
// 			src, err := file.Open()
// 			if err != nil {
// 				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			}
// 			defer src.Close()

// 			// Create a new file on the server
// 			dst, err2 := os.Create("public/assets/images/courses/" + file.Filename)
// 			if err2 != nil {
// 				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			}
// 			defer dst.Close()

// 			// Copy the uploaded file content to the new file
// 			if _, err = io.Copy(dst, src); err != nil {
// 				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			}

// 			course.Title = title
// 			course.Description = description
// 			course.Price = price
// 			course.Image = file.Filename
// 			course.UpdatedBy = userLogin.UserID
// 			course.UpdatedAt = time.Now()

// 			if err := DB.Model(&course).Updates(&course).Error; err != nil {
// 				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 				return c.RedirectBack("")
// 			}
// 		}
// 	}

// 	return c.Redirect("/managements/mng-courses")

// }
