package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"strconv"
	"time"
)

func GetNewsPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetHiringNewsPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward

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
	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	if err := DB.Find(&districts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	if err := DB.Find(&wards).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	return c.Render("pages/news/index", fiber.Map{
		"This":           9,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"Years":          years,
		"User":           userLogin,
		"provinces":      provinces,
		"districts":      districts,
		"wards":          wards,
		"Ctx":            c,
	}, "layouts/main")
}

func APIGetNews(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiring")
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

	var hiringNews []models.HiringNews
	var totalItems int64

	DB := initializers.DB

	// Lấy dữ liệu với phân trang
	if err := DB.Model(&models.HiringNews{}).
		Joins("User").
		Joins("Province").Joins("District").Joins("Ward").
		Where("hiring_news.deleted = ?", false).
		Where("User.deleted = ?", false).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&hiringNews).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	hiringNewsIDs := make([]int, len(hiringNews))
	for i, news := range hiringNews {
		hiringNewsIDs[i] = news.HiringNewsID
	}

	var userHiringNewsList []models.UserHiringNews
	if err := DB.Where("hiring_news_id IN ?", hiringNewsIDs).
		Where("status IN ?", []int{0, 1}).
		Where("nhaccong_id", userLogin.UserID).
		Find(&userHiringNewsList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userHiringNewsMap := make(map[int][]models.UserHiringNews)
	for _, userHiringNews := range userHiringNewsList {
		userHiringNewsMap[userHiringNews.HiringNewsID] = append(userHiringNewsMap[userHiringNews.HiringNewsID], userHiringNews)
	}

	for i := range hiringNews {
		if applicants, found := userHiringNewsMap[hiringNews[i].HiringNewsID]; found {
			hiringNews[i].Applicants = applicants
		}
	}

	for i := range hiringNews {
		if len(hiringNews[i].Applicants) > 0 {
			hiringNews[i].ApplicantStatus = hiringNews[i].Applicants[0].Status
		} else {
			hiringNews[i].ApplicantStatus = 4
		}
	}

	var contracts []models.Contract
	if err := DB.Where("nhaccong_id", userLogin.UserID).
		Where("deleted", false).Find(&contracts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	for i := range hiringNews {
		for _, contract := range contracts {
			if hiringNews[i].Date.Equal(contract.Date) {
				hiringNews[i].ApplicantStatus = 5
				break
			}
		}
	}

	// Tính tổng số phần tử
	if err := DB.Model(&models.HiringNews{}).Joins("User").
		Where("hiring_news.deleted = ?", false).
		Where("User.deleted = ?", false).
		Count(&totalItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":       hiringNews,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}

func PostNewsApply(c *fiber.Ctx) error {
	type Request struct {
		ID uint `json:"id"`
	}

	var req Request

	// Parse dữ liệu JSON từ request body
	if err := c.BodyParser(&req); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	hiringNewsID := req.ID
	DB := initializers.DB
	var userHiringNewsList []models.UserHiringNews
	userLogin := GetSessionUser(c)

	if err := DB.Where("hiring_news_id", hiringNewsID).Where("nhaccong_id", userLogin.UserID).
		Find(&userHiringNewsList).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	if len(userHiringNewsList) > 0 {
		userHiringNewsList[0].Status = 0

		if err := DB.Save(&userHiringNewsList[0]).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi khi cập nhật",
			})
		}
	} else {
		newUserHiringNews := models.UserHiringNews{
			HiringNewsID: int(hiringNewsID),
			NhaccongID:   userLogin.UserID,
			Status:       0,
			ApplyAt:      time.Now(),
		}

		if err := DB.Create(&newUserHiringNews).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi khi tạo yêu cầu ứng tuyển",
			})
		}
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Gửi yêu cầu ứng tuyển thành công",
	})
}

func PostNewsCancelApply(c *fiber.Ctx) error {
	type Request struct {
		ID uint `json:"id"`
	}

	var req Request

	// Parse dữ liệu JSON từ request body
	if err := c.BodyParser(&req); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	hiringNewsID := req.ID
	DB := initializers.DB
	var userHiringNewsList []models.UserHiringNews
	userLogin := GetSessionUser(c)

	if err := DB.Where("hiring_news_id", hiringNewsID).Where("nhaccong_id", userLogin.UserID).
		Find(&userHiringNewsList).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	if len(userHiringNewsList) > 0 {
		userHiringNewsList[0].Status = 2
		if err := DB.Save(&userHiringNewsList[0]).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi",
			})
		}
	} else {
		newUserHiringNews := models.UserHiringNews{
			HiringNewsID: int(hiringNewsID),
			NhaccongID:   userLogin.UserID,
			Status:       2,
			ApplyAt:      time.Now(),
		}

		if err := DB.Create(&newUserHiringNews).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi",
			})
		}
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Thu hồi yêu cầu ứng tuyển thành công",
	})
}
