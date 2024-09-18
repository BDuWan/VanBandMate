package controllers

import (
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"strconv"
	"time"

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

func GetEditContractPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get UserInfo")
	DB := initializers.DB
	var roles []models.Role
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward
	userLogin := GetSessionUser(c)
	user := new(models.User)
	userId := c.Params("id")
	if err := DB.Model(&models.User{}).Joins("Role").Joins(
		"Province").Joins("District").Joins("Ward").Where(
		"user_id", userId).First(user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Cannot get user")
	}
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	if err := DB.Where("deleted = ?", false).Find(&roles).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
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

	return c.Render("pages/management/mng-user/edit-user", fiber.Map{
		"This":        1,
		"Permissions": permissions,
		"User":        userLogin,
		"UserEdit":    user,
		"roles":       roles,
		"provinces":   provinces,
		"districts":   districts,
		"wards":       wards,
		"Ctx":         c,
	}, "layouts/main")
}

func APIPutEditContract(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "APIPutEditUser")
	var updateInfoForm structs.AdminUpdateInfoForm
	var account models.User
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	if err := c.BodyParser(&updateInfoForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Format User Fail")
	}
	validator := ValidatorAdminUpdateInfoInput(updateInfoForm)
	if validator != "ok" {
		return c.JSON(validator)
	}
	layout := "02/01/2006"
	date, err := time.Parse(layout, updateInfoForm.DateOfBirth)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: error parsing date of birth")
	}

	if err := DB.Where("user_id", updateInfoForm.UserID).First(&account).Error; err == nil {
		if account.Email != updateInfoForm.Email {
			if err := DB.Where("email", updateInfoForm.Email).First(&models.User{}).Error; err == nil {
				return c.JSON("Email đã được sử dụng, vui lòng chọn 1 email khác")
			}
			account.Verify = false
		}
		account.Gender = updateInfoForm.Gender
		account.FirstName = updateInfoForm.FirstName
		account.LastName = updateInfoForm.LastName
		account.Email = updateInfoForm.Email
		account.PhoneNumber = updateInfoForm.PhoneNumber
		account.LinkFacebook = updateInfoForm.LinkFacebook
		account.ProvinceCode = updateInfoForm.ProvinceCode
		account.DistrictCode = updateInfoForm.DistrictCode
		account.WardCode = updateInfoForm.WardCode
		account.AddressDetail = updateInfoForm.AddressDetail
		account.RoleID = updateInfoForm.RoleID
		account.Verify = updateInfoForm.Verify
		account.DateOfBirth = date
		account.UpdatedAt = time.Now()
		account.UpdatedBy = userLogin.UserID

		imageName := "avatar" + strconv.Itoa(account.UserID) + ".jpg"

		if updateInfoForm.Image != "" {
			path := "public/assets/img/avatar/"
			saveImageResult := SaveImage(updateInfoForm.Image, path, imageName)
			if saveImageResult != "ok" {
				return c.JSON(saveImageResult)
			}
			account.Image = imageName
		}
		if err := DB.Save(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Can not create account")
			return c.JSON("Không thể cập nhật thông tin")
		}

		return c.JSON("Success")
	}
	return c.JSON("Đã xảy ra lỗi khi lấy thông tin người dùng")
}

func APIDeleteContract(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIDeleteUserID")
	var user models.User
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	userId := c.Params("id")

	idUser, err := strconv.Atoi(userId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

	}
	if idUser == 1 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: User=1")

		return c.JSON("Tài khoản này không thể xóa")
	}

	if err := DB.Where("user_id", userId).Where("deleted", false).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]:" + err.Error())
		return c.JSON("Không tìm thấy tài khoản")
	}

	user.Deleted = true
	user.DeletedBy = userLogin.UserID
	user.DeletedAt = time.Now()

	if err := DB.Model(&user).Updates(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]:" + err.Error())
		return c.JSON("Can not delete this role.")
	}

	return c.JSON("success")
}
