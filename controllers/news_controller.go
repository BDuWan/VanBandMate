package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"lms/initializers"
	"lms/models"
	"lms/structs"
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

//func APIGetNews(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiring")
//	userLogin := GetSessionUser(c)
//	// Lấy số trang và số phần tử mỗi trang từ query params, với giá trị mặc định nếu không có
//	page, err := strconv.Atoi(c.Query("page", "1"))
//	if err != nil || page < 1 {
//		page = 1
//	}
//	limit, err := strconv.Atoi(c.Query("itemsPerPage", "1"))
//	if err != nil || limit < 1 {
//		limit = 5
//	}
//
//	offset := (page - 1) * limit
//
//	var hiringNews []models.HiringNews
//	var totalItems int64
//
//	DB := initializers.DB
//
//	// Lấy dữ liệu với phân trang
//	if err := DB.Model(&models.HiringNews{}).
//		Joins("User").
//		Joins("Province").Joins("District").Joins("Ward").
//		Where("hiring_news.deleted = ?", false).
//		Where("User.deleted = ?", false).
//		Order("created_at DESC").
//		Offset(offset).Limit(limit).Find(&hiringNews).Error; err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"error": err.Error(),
//		})
//	}
//
//	hiringNewsIDs := make([]int, len(hiringNews))
//	for i, news := range hiringNews {
//		hiringNewsIDs[i] = news.HiringNewsID
//	}
//
//	var userHiringNewsList []models.UserHiringNews
//	if err := DB.Where("hiring_news_id IN ?", hiringNewsIDs).
//		Where("status IN ?", []int{0, 1}).
//		Where("nhaccong_id", userLogin.UserID).
//		Find(&userHiringNewsList).Error; err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"error": err.Error(),
//		})
//	}
//
//	//Kiểm tra trạng thái ứng tuyển của bản thân với từng tin tuyển dụng
//	userHiringNewsMap := make(map[int][]models.UserHiringNews)
//	for _, userHiringNews := range userHiringNewsList {
//		userHiringNewsMap[userHiringNews.HiringNewsID] = append(userHiringNewsMap[userHiringNews.HiringNewsID], userHiringNews)
//	}
//
//	for i := range hiringNews {
//		if applicants, found := userHiringNewsMap[hiringNews[i].HiringNewsID]; found {
//			hiringNews[i].Applicants = applicants
//		}
//	}
//
//	for i := range hiringNews {
//		if len(hiringNews[i].Applicants) > 0 {
//			hiringNews[i].ApplicantStatus = hiringNews[i].Applicants[0].Status
//		} else {
//			hiringNews[i].ApplicantStatus = 4
//		}
//	}
//
//	//Kiểm tra tin trùng thời gian
//	var contracts []models.Contract
//	if err := DB.Where("nhaccong_id", userLogin.UserID).
//		Where("deleted", false).Find(&contracts).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//		return c.JSON("Đã xảy ra lỗi")
//	}
//
//	for i := range hiringNews {
//		for _, contract := range contracts {
//			if hiringNews[i].Date.Equal(contract.Date) {
//				hiringNews[i].ApplicantStatus = 5
//				break
//			}
//		}
//	}
//
//	for i := range hiringNews {
//		if len(hiringNews[i].Applicants) > 0 {
//			if hiringNews[i].Applicants[0].Status == 1 {
//				hiringNews[i].ApplicantStatus = 1
//			}
//		}
//	}
//
//	// Tính tổng số phần tử
//	if err := DB.Model(&models.HiringNews{}).Joins("User").
//		Where("hiring_news.deleted = ?", false).
//		Where("User.deleted = ?", false).
//		Count(&totalItems).Error; err != nil {
//		return c.Status(500).JSON(fiber.Map{
//			"error": err.Error(),
//		})
//	}
//
//	return c.JSON(fiber.Map{
//		"data":       hiringNews,
//		"user_id":    userLogin.UserID,
//		"page":       page,
//		"totalItems": totalItems,
//	})
//}

func APIPostNewsFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostNewsFilter")

	form := new(structs.FormFilterNews)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	page := form.Page
	limit := form.ItemsPerPage

	offset := (page - 1) * limit

	DB := initializers.DB
	userLogin := GetSessionUser(c)
	query := DB.Model(&models.HiringNews{}).
		Joins("User").Joins("Province").Joins("District").Joins("Ward").
		Where("hiring_news.deleted", false).
		Where("User.deleted", false).
		Where("hiring_news.date > ?", time.Now())

	if form.HiringEnough != 2 {
		query = query.Where("hiring_news.hiring_enough", form.HiringEnough)
	}

	if form.Year != 0 {
		query = query.Where("YEAR(hiring_news.date) = ?", form.Year)
	}

	if form.Month != 0 {
		query = query.Where("MONTH(hiring_news.date) = ?", form.Month)
	}

	if form.TimeCreate == 1 {
		query = query.Where("DATE(hiring_news.created_at) = CURDATE()")
	} else if form.TimeCreate == 2 {
		query = query.Where("MONTH(hiring_news.created_at) = MONTH(CURDATE()) AND YEAR(hiring_news.created_at) = YEAR(CURDATE())")
	} else if form.TimeCreate == 3 {
		query = query.Where("YEAR(hiring_news.created_at) = YEAR(CURDATE())")
	}

	if form.ProvinceCode != "0" && form.ProvinceCode != "" {
		query = query.Where("hiring_news.province_code", form.ProvinceCode)
	}

	if form.DistrictCode != "0" && form.DistrictCode != "" {
		query = query.Where("hiring_news.district_code", form.DistrictCode)
	}

	if form.WardCode != "0" && form.WardCode != "" {
		query = query.Where("hiring_news.ward_code", form.WardCode)
	}

	if form.Order == 0 {
		query = query.Order("created_at DESC")
	} else if form.Order == 1 {
		query = query.Order("update_at DESC")
	} else if form.Order == 2 {
		query = query.Order("date ASC")
	} else if form.Order == 3 {
		query = query.Order("price DESC")
	}

	var totalItems int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&totalItems).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var hiringNews []models.HiringNews
	if err := query.Offset(offset).Limit(limit).Find(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
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

	//Kiểm tra trạng thái ứng tuyển của bản thân với từng tin tuyển dụng
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

	//Kiểm tra tin trùng thời gian
	var contracts []models.Contract
	if err := DB.Where("nhaccong_id", userLogin.UserID).
		Where("deleted", false).Find(&contracts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	for i := range hiringNews {
		for _, contract := range contracts {
			if hiringNews[i].Date.Equal(contract.Date) {
				if hiringNews[i].ApplicantStatus != 1 {
					hiringNews[i].ApplicantStatus = 5
					break
				}
			}
		}
	}

	//for i := range hiringNews {
	//	if len(hiringNews[i].Applicants) > 0 {
	//		if hiringNews[i].Applicants[0].Status == 1 {
	//			hiringNews[i].ApplicantStatus = 1
	//		}
	//	}
	//}

	return c.JSON(fiber.Map{
		"message":    "success",
		"data":       hiringNews,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}

func PostNewsApply(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostNewsApply")
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
	var hiringNews models.HiringNews

	if err := DB.Where("hiring_news_id", hiringNewsID).
		Where("deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Tin tuyển dụng đã bị xóa",
		})
	}

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
			Date:         hiringNews.Date,
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
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostNewsCancelApply")
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
	//var userHiringNewsList []models.UserHiringNews
	userLogin := GetSessionUser(c)

	if err := DB.Model(&models.UserHiringNews{}).
		Where("hiring_news_id", hiringNewsID).Where("nhaccong_id", userLogin.UserID).
		Update("status", 2).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Thu hồi yêu cầu ứng tuyển thành công",
	})
}
