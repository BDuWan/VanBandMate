package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"strconv"
	"strings"
	"time"
)

func GetHiringPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetHiringPage")
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
	return c.Render("pages/hiring/index", fiber.Map{
		"This":           7,
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

func APIGetHiring(c *fiber.Ctx) error {
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
	if err := DB.Model(&models.HiringNews{}).Joins("User").
		Joins("Province").Joins("District").Joins("Ward").
		Where("hiring_news.deleted = ?", false).
		Where("User.deleted = ?", false).
		Where("hiring_news.chuloadai_id", userLogin.UserID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&hiringNews).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Tính tổng số phần tử
	if err := DB.Model(&models.HiringNews{}).Joins("User").
		Where("hiring_news.deleted = ?", false).
		Where("User.deleted = ?", false).
		Where("hiring_news.chuloadai_id", userLogin.UserID).
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

func APIPostHiringFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostHiringFilter")

	form := new(structs.FormFilter)
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
		Where("hiring_news.deleted", false).Where("User.deleted", false)

	if form.Employer == 1 {
		query = query.Where("hiring_news.chuloadai_id", userLogin.UserID)
	}

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

	if form.Order == 0 {
		query = query.Order("created_at DESC")
	} else if form.Order == 1 {
		query = query.Order("update_at DESC")
	} else if form.Order == 2 {
		query = query.Order("date ASC")
	} else if form.Order == 3 {
		query = query.Order("price DESC")
	} else if form.Order == 4 {
		query = query.Order("ward_code ASC")
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

	return c.JSON(fiber.Map{
		"message":    "success",
		"data":       hiringNews,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}

func APIPostHiringCreate(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostCreateRole")
	var hiringNews models.HiringNews
	userLogin := GetSessionUser(c)

	DB := initializers.DB
	var hiringNewsForm structs.HiringNewsForm
	if err := c.BodyParser(&hiringNewsForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Dữ liệu không hợp lệ")
	}
	layout := "02/01/2006"
	date, _ := time.Parse(layout, hiringNewsForm.Date)

	valid := ValidatorHiringNewsInput(hiringNewsForm)
	if valid != "ok" {
		return c.JSON(valid)
	}
	price, _ := strconv.Atoi(hiringNewsForm.Price)

	hiringNews.ChuloadaiID = userLogin.UserID
	hiringNews.ProvinceCode = hiringNewsForm.Province
	hiringNews.DistrictCode = hiringNewsForm.District
	hiringNews.WardCode = hiringNewsForm.Ward
	hiringNews.Date = date
	hiringNews.Describe = hiringNewsForm.Describe
	hiringNews.Price = price
	hiringNews.AddressDetail = hiringNewsForm.Address
	hiringNews.Deleted = false
	hiringNews.CreatedBy = userLogin.UserID
	hiringNews.DeletedAt = time.Now()
	hiringNews.CreatedAt = time.Now()

	if err := DB.Create(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Không thể tạo tin")
	}
	return c.JSON("success")
}

func APIPutHiringUpdate(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPutHiringUpdate")
	var hiringNews models.HiringNews
	hiringNewsId := c.Params("id")
	userLogin := GetSessionUser(c)

	DB := initializers.DB
	var hiringNewsForm structs.HiringNewsForm
	if err := c.BodyParser(&hiringNewsForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Dữ liệu không hợp lệ")
	}
	if err := DB.Where("hiring_news_id = ?", hiringNewsId).Where("deleted", false).First(&hiringNews).Error; err != nil {

		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Tin đã bị xóa")
	}
	if hiringNews.ChuloadaiID != userLogin.UserID {
		return c.JSON("Bạn không thể chỉnh sửa tin của người khác")
	}
	layout := "02/01/2006"
	date, _ := time.Parse(layout, hiringNewsForm.Date)

	valid := ValidatorHiringNewsInput(hiringNewsForm)
	if valid != "ok" {
		return c.JSON(valid)
	}
	price, _ := strconv.Atoi(hiringNewsForm.Price)

	hiringNews.ProvinceCode = hiringNewsForm.Province
	hiringNews.DistrictCode = hiringNewsForm.District
	hiringNews.WardCode = hiringNewsForm.Ward
	hiringNews.Date = date
	hiringNews.Describe = hiringNewsForm.Describe
	hiringNews.Price = price
	hiringNews.AddressDetail = hiringNewsForm.Address
	hiringNews.UpdatedAt = time.Now()
	hiringNews.UpdatedBy = userLogin.UserID

	if err := DB.Model(&hiringNews).Updates(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Đã xảy ra lỗi trong quá trình lưu dữ liệu ")
	}
	return c.JSON("success")
}

func APIGetHiringID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringNewsID")
	var hiringNews models.HiringNews
	hiringNewsId := c.Params("id")
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if err := DB.Model(&models.HiringNews{}).
		Where("hiring_news_id", hiringNewsId).Where("deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}
	if userLogin.UserID != hiringNews.ChuloadaiID {
		return c.JSON(fiber.Map{
			"message": "Bạn chỉ được phép chỉnh sửa những tin do chính bạn tạo",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    hiringNews,
	})
}

func APIGetHiringDetailID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringNewsID")
	var hiringNews models.HiringNews
	hiringNewsId := c.Params("id")
	DB := initializers.DB

	if err := DB.Model(&models.HiringNews{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("User").Joins("User.Role").
		Where("hiring_news_id", hiringNewsId).Where("hiring_news.deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    hiringNews,
	})
}

func APIGetHiringListApply(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringListApply")
	var userHiringNews []models.UserHiringNews
	var hiringNews models.HiringNews
	hiringNewsId := c.Params("id")
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	if err := DB.Where("hiring_news_id", hiringNewsId).
		Where("hiring_news.deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	if hiringNews.ChuloadaiID != userLogin.UserID {
		return c.JSON(fiber.Map{
			"message": "Bạn chỉ được xem danh sách ứng tuyển của tin mà bạn tạo",
		})
	}

	if err := DB.Model(&models.UserHiringNews{}).Joins("User").
		Joins("User.Province").Joins("User.District").Joins("User.Ward").
		Where("hiring_news_id", hiringNewsId).
		Where("user_hiring_news.status IN ?", []int{0, 1}).
		Find(&userHiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "success",
		"data":          userHiringNews,
		"hiring_enough": hiringNews.HiringEnough,
	})
}

func APIPostSaveApply(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostSaveApply")
	var hiringNews models.HiringNews
	userLogin := GetSessionUser(c)

	DB := initializers.DB
	var formSaveApply structs.FormSaveApply

	if err := c.BodyParser(&formSaveApply); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Dữ liệu không hợp lệ")
	}

	if err := DB.Where("hiring_news_id = ?", formSaveApply.HiringNewsID).
		Where("deleted", false).First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Tin đã bị xóa")
	}

	if formSaveApply.HiringEnough != hiringNews.HiringEnough {
		hiringNews.HiringEnough = formSaveApply.HiringEnough

		if err := DB.Save(&hiringNews).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Không thể cập nhật trạng thái tuyển đủ người")
		}
	}

	if len(formSaveApply.SelectedItems) > 0 {
		var contracts []models.Contract
		var dupContracts []models.Contract
		var nhaccongIDs []int
		for _, item := range formSaveApply.SelectedItems {
			nhaccongIDs = append(nhaccongIDs, item.NhaccongID)
		}

		//Kiểm tra xem có ai trùng ngày không
		if err := DB.Model(&models.Contract{}).
			Where("nhaccong_id IN ?", nhaccongIDs).
			Where("DATE(date) = ?", hiringNews.Date.Format("2006-01-02")).
			Where("deleted", false).
			Find(&dupContracts).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi")
		}

		numberDupContracts := len(dupContracts)

		if numberDupContracts > 0 {
			return c.JSON("Nhạc công trong danh sách bạn chọn đã tồn tại hợp đồng trùng với ngày trên tin tuyển dụng")
		}

		// Lặp qua các SelectedItems và tạo hợp đồng
		for _, item := range formSaveApply.SelectedItems {
			contract := models.Contract{
				ChuloadaiID:   hiringNews.ChuloadaiID,
				NhaccongID:    item.NhaccongID,
				ProvinceCode:  hiringNews.ProvinceCode,
				DistrictCode:  hiringNews.DistrictCode,
				WardCode:      hiringNews.WardCode,
				AddressDetail: hiringNews.AddressDetail,
				Status:        1,
				Price:         hiringNews.Price,
				Date:          hiringNews.Date,
				Deleted:       false,
				CreatedBy:     userLogin.UserID,
				CreatedAt:     time.Now(),
			}
			contracts = append(contracts, contract)
		}

		if err := DB.Create(&contracts).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình tạo hợp đồng")
		}

		var userHiringNewsIDs []int
		for _, item := range formSaveApply.SelectedItems {
			userHiringNewsIDs = append(userHiringNewsIDs, item.UserHiringNewsID)
		}

		var userApplyIDs []int
		for _, item := range formSaveApply.SelectedItems {
			userApplyIDs = append(userHiringNewsIDs, item.NhaccongID)
		}

		////Cập nhật các yêu cầu ứng tuyển liên quan
		//var userHiringNewsList []models.UserHiringNews
		//if err := DB.Model(&models.UserHiringNews{}).
		//	Joins("JOIN hiring_news ON user_hiring_news.hiring_news_id = hiring_news.hiring_news_id").
		//	Select("user_hiring_news.*, hiring_news.date").
		//	Where("user_hiring_news.nhaccong_id IN ?", userApplyIDs).
		//	Where("DATE(hiring_news.date) = ?", hiringNews.Date.Format("2006-01-02")).
		//	Find(&userHiringNewsList).Error; err != nil {
		//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		//	return c.JSON("Đã xảy ra lỗi trong quá trình lấy dữ liệu")
		//}
		//
		//// In ra danh sách các bản ghi lấy được
		//fmt.Println("Danh sách user_hiring_news lấy được:")
		//for _, item := range userHiringNewsList {
		//	fmt.Printf("ID: %d, Date: %s, Status: %d\n", item.UserHiringNewsID, item.Status)
		//}

		if err := DB.Model(&models.UserHiringNews{}).
			Where("nhaccong_id IN ?", userApplyIDs).
			Where("DATE(date) = ?", hiringNews.Date.Format("2006-01-02")).
			Updates(models.UserHiringNews{Status: 3}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}

		if err := DB.Model(&models.UserHiringNews{}).
			Where("user_hiring_news_id IN ?", userHiringNewsIDs).
			Updates(models.UserHiringNews{Status: 1}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}
	}

	return c.JSON("success")
}

func ValidatorHiringNewsInput(hiringNewsForm structs.HiringNewsForm) string {
	if hiringNewsForm.Province == "0" {
		return "Bạn chưa chọn tỉnh/thành phố"
	}
	if hiringNewsForm.District == "0" {
		return "Bạn chưa chọn quận/huyện"
	}
	if hiringNewsForm.Ward == "0" {
		return "Bạn chưa chọn xã/phường/thị trấn"
	}
	if hiringNewsForm.Date == "" {
		return "Bạn chưa nhập thời gian"
	}

	if strings.TrimSpace(hiringNewsForm.Address) == "" {
		return "Bạn chưa nhập địa chỉ chi tiết"
	}

	layout := "02/01/2006"
	date, err := time.Parse(layout, hiringNewsForm.Date)
	if err != nil {
		return "Định dạng thời gian không hợp lệ"
	}

	if time.Now().After(date) {
		return "Không thể tạo tin tuyển dụng của ngày trong quá khứ"
	}

	price, err := strconv.Atoi(hiringNewsForm.Price)
	if err != nil {
		return "Giá tiền không hợp lệ"
	}

	if price <= 0 {
		return "Giá tiền phải là số nguyên dương và lớn hơn 0"
	}

	return "ok"
}
