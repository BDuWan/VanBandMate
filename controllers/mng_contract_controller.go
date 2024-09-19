package controllers

import (
	"time"
	"vanbandmate/initializers"
	"vanbandmate/models"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
)

func GetMngContractPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetMngUserPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	var provinces []models.Province
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}

	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
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

	return c.Render("pages/management/mng-contract/index", fiber.Map{
		"This":           3,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"Years":          years,
		"Provinces":      provinces,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

func APIPostMngContractFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostMngContractFilter")
	var contracts []models.Contract
	DB := initializers.DB
	form := new(struct {
		Status   int    `json:"status"`
		Province string `json:"province_code"`
		Year     int    `json:"year"`
		Month    int    `json:"month"`
	})
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM1]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	query := DB.Model(&models.Contract{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("ChuLoaDai").Joins("NhacCong")

	if form.Status != 4 {
		query = query.Where("contracts.status", form.Status)
	}

	if form.Province != "0" {
		query = query.Where("contracts.province_code", form.Province)
	}

	if form.Year != 0 {
		query = query.Where("YEAR(contracts.date) = ?", form.Year)
	}

	if form.Month != 0 {
		query = query.Where("MONTH(contracts.date) = ?", form.Month)
	}

	if err := query.Find(&contracts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM2]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	return c.JSON(fiber.Map{
		"data": contracts,
	})
}

func APIGetMngContractDetailID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetHiringNewsID")
	var contract models.Contract
	contractId := c.Params("id")
	DB := initializers.DB

	if err := DB.Model(&models.Contract{}).
		Joins("Province").Joins("District").Joins("Ward").
		Joins("ChuLoaDai").Joins("NhacCong").
		Where("contract_id", contractId).
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

func APIPutMngContractRestoreID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "APIPutMngContractRestoreID")

	contractID := c.Params("id")
	DB := initializers.DB

	var contract models.Contract
	var numberDupContract int64

	if err := DB.Model(&models.Contract{}).
		Where("contract_id = ?", contractID).
		First(&contract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	if err := DB.Model(&models.Contract{}).
		Where("nhaccong_id = ?", contract.NhaccongID).
		Where("DATE(date) = ?", contract.Date.Format("2006-01-02")).
		Where("deleted = ?", false).
		Count(&numberDupContract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	if numberDupContract > 0 {
		return c.JSON("Nhạc công đã có hợp đồng khác trong ngày này")
	}
	if contract.Date.After(time.Now()) {
		if err := DB.Model(&contract).Updates(map[string]interface{}{
			"deleted": 0,
			"status":  1,
		}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}
	} else {
		if err := DB.Model(&contract).Updates(map[string]interface{}{
			"deleted": 0,
			"status":  0,
		}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}
	}

	if err := DB.Model(&models.UserHiringNews{}).
		Where("DATE(date) = ?", contract.Date.Format("2006-01-02")).
		Where("nhaccong_id", contract.NhaccongID).
		Where("status = ?", 0).
		Select("status").
		Updates(&models.UserHiringNews{Status: 3}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	// Trả về phản hồi thành công
	return c.JSON("success")
}

func APIPutMngContractDeleteID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "APIPutMngContractDeleteID")

	contractID := c.Params("id")
	DB := initializers.DB

	var contract models.Contract

	if err := DB.Model(&models.Contract{}).
		Where("contract_id = ?", contractID).
		First(&contract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	if contract.Date.After(time.Now()) {
		if err := DB.Model(&contract).Updates(map[string]interface{}{
			"deleted": 1,
			"status":  3,
		}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}
	} else {
		if err := DB.Model(&contract).Updates(map[string]interface{}{
			"deleted": 1,
			"status":  3,
		}).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi trong quá trình cập nhật dữ liệu")
		}
	}

	if err := DB.Model(&models.UserHiringNews{}).
		Where("DATE(date) = ?", contract.Date.Format("2006-01-02")).
		Where("nhaccong_id", contract.NhaccongID).
		Where("status = ?", 3).
		Select("status").
		Updates(&models.UserHiringNews{Status: 0}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	// Trả về phản hồi thành công
	return c.JSON("success")
}
