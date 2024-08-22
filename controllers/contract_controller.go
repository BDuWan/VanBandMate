package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"strconv"
	"time"
)

func GetMyContractPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetMyContractPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	var years []int
	if err := DB.Model(&models.Contract{}).
		Where("date IS NOT NULL").
		Where("deleted", false).
		Select("DISTINCT EXTRACT(YEAR FROM date) AS year").
		Order("year ASC").
		Pluck("year", &years).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("error")
	}

	return c.Render("pages/contract/index", fiber.Map{
		"This":           6,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"Years":          years,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

func APIGetContract(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetContract")
	userLogin := GetSessionUser(c)
	// Lấy số trang và số phần tử mỗi trang từ query params, với giá trị mặc định nếu không có
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("itemsPerPage", "1"))
	if err != nil || limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit

	var contracts []models.Contract
	var totalItems int64

	DB := initializers.DB

	// Lấy dữ liệu với phân trang
	if err := DB.Model(&models.Contract{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("ChuLoaDai").Joins("NhacCong").
		Where("contracts.deleted = ?", false).
		Where("ChuLoaDai.deleted = ?", false).
		Where("NhacCong.deleted = ?", false).
		Order("status DESC").
		Offset(offset).Limit(limit).Find(&contracts).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Tính tổng số phần tử
	if err := DB.Model(&models.Contract{}).Joins("ChuLoaDai").Joins("NhacCong").
		Where("contracts.deleted = ?", false).
		Where("ChuLoaDai.deleted = ?", false).
		Where("NhacCong.deleted = ?", false).
		Count(&totalItems).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":       contracts,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}
