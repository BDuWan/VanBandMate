package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"time"
	"vanbandmate/initializers"
	"vanbandmate/models"
	"vanbandmate/structs"
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

//func APIGetContract(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetContract")
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
//	var contracts []models.Contract
//	var totalItems int64
//
//	DB := initializers.DB
//
//	// Lấy dữ liệu với phân trang
//	if err := DB.Model(&models.Contract{}).
//		Joins("Province").Joins("District").Joins("Ward").
//		Joins("ChuLoaDai").Joins("NhacCong").
//		Where("contracts.chuloadai_id = ? OR contracts.nhaccong_id = ?", userLogin.UserID, userLogin.UserID).
//		Where("contracts.deleted = ?", false).
//		Where("ChuLoaDai.deleted = ?", false).
//		Where("NhacCong.deleted = ?", false).
//		Order("status DESC").
//		Offset(offset).Limit(limit).Find(&contracts).Error; err != nil {
//		return c.JSON(fiber.Map{
//			"error": err.Error(),
//		})
//	}
//
//	// Tính tổng số phần tử
//	if err := DB.Model(&models.Contract{}).
//		Joins("ChuLoaDai").Joins("NhacCong").
//		Where("contracts.chuloadai_id = ? OR contracts.nhaccong_id = ?", userLogin.UserID, userLogin.UserID).
//		Where("contracts.deleted = ?", false).
//		Where("ChuLoaDai.deleted = ?", false).
//		Where("NhacCong.deleted = ?", false).
//		Count(&totalItems).Error; err != nil {
//		return c.JSON(fiber.Map{
//			"error": err.Error(),
//		})
//	}
//
//	return c.JSON(fiber.Map{
//		"data":       contracts,
//		"user_id":    userLogin.UserID,
//		"page":       page,
//		"totalItems": totalItems,
//	})
//}

func APIGetContractDetailID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringNewsID")
	var contract models.Contract
	contractId := c.Params("id")
	DB := initializers.DB

	if err := DB.Model(&models.Contract{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("ChuLoaDai").Joins("NhacCong").
		Where("contract_id", contractId).Where("contracts.deleted", false).
		First(&contract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    contract,
	})
}

func APIPostContractFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostContractFilter")

	form := new(structs.FormFilterContract)
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
	query := DB.Model(&models.Contract{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("ChuLoaDai").Joins("NhacCong").
		Where("contracts.chuloadai_id = ? OR contracts.nhaccong_id = ?", userLogin.UserID, userLogin.UserID).
		Where("contracts.deleted = ?", false).
		Where("ChuLoaDai.deleted = ?", false).
		Where("NhacCong.deleted = ?", false)

	if form.Status != 3 {
		query = query.Where("contracts.status", form.Status)
	}

	if form.Year != 0 {
		query = query.Where("YEAR(contracts.date) = ?", form.Year)
	}

	if form.Month != 0 {
		query = query.Where("MONTH(contracts.date) = ?", form.Month)
	}

	if form.TimeCreate == 1 {
		query = query.Where("DATE(contracts.created_at) = CURDATE()")
	} else if form.TimeCreate == 2 {
		query = query.Where("MONTH(contracts.created_at) = MONTH(CURDATE()) AND YEAR(contracts.created_at) = YEAR(CURDATE())")
	} else if form.TimeCreate == 3 {
		query = query.Where("YEAR(contracts.created_at) = YEAR(CURDATE())")
	}

	if form.Order == 0 {
		query = query.Order("status DESC")
	} else if form.Order == 1 {
		query = query.Order("date ASC")
	} else if form.Order == 2 {
		query = query.Order("created_at DESC")
	} else if form.Order == 3 {
		query = query.Order("price ASC")
	} else if form.Order == 4 {
		query = query.Order("ward_code ASC")
	}

	var totalItems int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&totalItems).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var contracts []models.Contract
	if err := query.Offset(offset).Limit(limit).Find(&contracts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	return c.JSON(fiber.Map{
		"message":    "success",
		"data":       contracts,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}

func PostContractRequestDelete(c *fiber.Ctx) error {
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

	contractID := req.ID
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	if err := DB.Model(&models.Contract{}).
		Where("contract_id = ?", contractID).
		Where("deleted = ?", false).
		Updates(map[string]interface{}{
			"status":            2,
			"request_delete_by": userLogin.UserID,
		}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không tìm thấy hợp đồng",
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Gửi yêu cầu hủy thành công",
	})
}

func PostContractCancelDelete(c *fiber.Ctx) error {
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

	contractID := req.ID
	DB := initializers.DB

	if err := DB.Model(&models.Contract{}).
		Where("contract_id = ?", contractID).
		Where("deleted = ?", false).
		Update("status", 1).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không tìm thấy hợp đồng",
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Thu hồi yêu cầu hủy thành công",
	})
}

func PostContractConfirmDelete(c *fiber.Ctx) error {
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

	contractID := req.ID
	DB := initializers.DB

	var contract models.Contract

	if err := DB.Model(&models.Contract{}).
		Where("contract_id = ?", contractID).
		Where("deleted = ?", false).
		First(&contract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không tìm thấy hợp đồng",
		})
	}

	if err := DB.Model(&contract).Updates(map[string]interface{}{
		"deleted": 1,
		"status":  3,
	}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
	}

	if contract.RequestDeleteBy == contract.ChuloadaiID {
		if err := DB.Model(&models.UserHiringNews{}).
			Where("DATE(date) = ?", contract.Date.Format("2006-01-02")).
			Where("nhaccong_id", contract.NhaccongID).
			Where("status = ?", 3).
			Select("status").
			Updates(&models.UserHiringNews{Status: 0}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi",
			})
		}
	}

	if err := DB.Model(&models.UserHiringNews{}).
		Where("DATE(date) = ?", contract.Date.Format("2006-01-02")).
		Where("nhaccong_id", contract.NhaccongID).
		Where("status", 1).
		Updates(models.UserHiringNews{Status: 2}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Hợp đồng đã được hủy thành công",
	})
}
