package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"strconv"
	"time"
)

func GetHiringPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetHiringPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	var years []int
	if err := DB.Model(&models.HiringNews{}).
		Where("date IS NOT NULL").
		Where("deleted", false).
		Select("DISTINCT EXTRACT(YEAR FROM date) AS year").
		Order("year ASC").
		Pluck("year", &years).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("error")
	}
	return c.Render("pages/hiring/index", fiber.Map{
		"This":           7,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"Years":          years,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

//	APIGetHiring func APIGetHiring(c *fiber.Ctx) error {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringNews")
//		var hiringNews []models.HiringNews
//		DB := initializers.DB
//
//		if err := DB.Model(&models.HiringNews{}).Joins(
//			"User").Joins("Province").Joins("District").Joins("Ward").Where(
//			"hiring_news.deleted", false).Find(&hiringNews).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//		}
//		userLogin := GetSessionUser(c)
//
//		return c.JSON(fiber.Map{
//			"hiringNews": hiringNews,
//			"User":       userLogin,
//		})
//	}
func APIGetHiring(c *fiber.Ctx) error {
	// Lấy số trang và số phần tử mỗi trang từ query params, với giá trị mặc định nếu không có
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("itemsPerPage", "1"))
	if err != nil || limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit

	var hiringNews []models.HiringNews
	var totalItems int64 // Sử dụng int64 cho COUNT

	DB := initializers.DB

	// Lấy dữ liệu với phân trang
	if err := DB.Model(&models.HiringNews{}).Joins("User").
		Joins("Province").Joins("District").Joins("Ward").
		Where("hiring_news.deleted = ?", false).
		Where("User.deleted = ?", false).
		Offset(offset).Limit(limit).Find(&hiringNews).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Tính tổng số phần tử
	if err := DB.Model(&models.HiringNews{}).Joins("User").Where(
		"hiring_news.deleted = ?", false).Where(
		"User.deleted = ?", false).Count(&totalItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":       hiringNews,
		"page":       page,
		"totalItems": totalItems,
	})
}
