package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"lms/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetProfile(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get Profile ID")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	return c.Render("pages/info/profile", fiber.Map{
		"Permissions": permissions,
		"User":        userLogin,
		"Ctx":         c,
	}, "layouts/main")
}

func GetProfileID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get Profile ID")
	userId := c.Params("id")
	DB := initializers.DB
	var user models.User
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	if err := DB.Model(&models.User{}).Joins("Role").Joins(
		"Province").Joins("District").Joins("Ward").Where(
		"user_id", userId).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Cannot get user")
	}

	return c.Render("pages/info/profile", fiber.Map{
		"Permissions": permissions,
		"User":        user,
		"Ctx":         c,
	}, "layouts/main")
}

func GetUserInfo(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get UserInfo")
	DB := initializers.DB
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward
	userLogin := GetSessionUser(c)
	user := new(models.User)
	if err := DB.Model(&models.User{}).Joins("Role").Joins(
		"Province").Joins("District").Joins("Ward").Where(
		"user_id", userLogin.UserID).First(user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Cannot get user")
	}
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

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

	return c.Render("pages/info/userinfo", fiber.Map{
		"Permissions": permissions,
		"User":        user,
		"provinces":   provinces,
		"districts":   districts,
		"wards":       wards,
		"Ctx":         c,
	}, "layouts/main")
}

func PutUserInfo(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "PutUserInfo")
	var updateInfoForm structs.UpdateInfoForm
	var account models.User
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	if err := c.BodyParser(&updateInfoForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Format User Fail")
	}
	validator := ValidatorUpdateInfoInput(updateInfoForm)
	if validator != "ok" {
		return c.JSON(validator)
	}
	layout := "02/01/2006"
	date, err := time.Parse(layout, updateInfoForm.DateOfBirth)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: error parsing date of birth")
	}

	if err := DB.Where("user_id", userLogin.UserID).First(&account).Error; err == nil {
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
		account.DateOfBirth = date
		account.UpdatedAt = time.Now()
		account.UpdatedBy = userLogin.UserID
		fmt.Println(account.Verify)

		imageName := "avatar" + strconv.Itoa(account.UserID) + ".jpg"

		if updateInfoForm.Image != "" {
			path := "public/assets/img/avatar/"
			saveImageResult := SaveImage(updateInfoForm.Image, path, imageName)
			if saveImageResult != "ok" {
				return c.JSON(saveImageResult)
			}
			account.Image = imageName
		} else {
			account.Image = "default.jpg"
		}
		if err := DB.Save(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Can not create account")
			return c.JSON("Không thể cập nhật thông tin")
		}

		return c.JSON("Success")
	}
	return c.JSON("Đã xảy ra lỗi khi lấy thông tin người dùng")
}

func ValidatorUpdateInfoInput(user structs.UpdateInfoForm) string {
	//name
	if strings.TrimSpace(user.LastName) == "" {
		return "Họ và tên đệm không được để trống"
	}
	if strings.TrimSpace(user.FirstName) == "" {
		return "Tên không được để trống"
	}
	regexName := "[0-9!@#$%^&*()_+?:;,./={}~]"
	regexN := regexp.MustCompile(regexName)
	if regexN.MatchString(user.FirstName) {
		return "Họ và tên đệm không được có số hoặc kí tự đặc biệt"
	}

	if regexN.MatchString(user.LastName) {
		return "Tên không được có số hoặc kí tự đặc biệt"
	}

	//sdt
	if strings.TrimSpace(user.PhoneNumber) == "" {
		return "Số điện thoại không được để trống"
	}
	regexPhone := "^[0-9]{10,}$"
	regexP := regexp.MustCompile(regexPhone)
	if !regexP.MatchString(user.PhoneNumber) {
		return "Vui lòng nhập số điện thoại hợp lệ"
	}

	//dob
	if user.DateOfBirth == "" {
		return "Vui lòng nhập ngày sinh"
	}

	layout := "02/01/2006"
	date, _ := time.Parse(layout, user.DateOfBirth)

	if time.Now().Year()-date.Year() < 18 {
		return "Cần đạt tối thiểu 18 tuổi"
	}

	//email
	if user.Email == "" {
		return "Email không được để trống"
	}
	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexE := regexp.MustCompile(regexEmail)
	if !regexE.MatchString(user.Email) {
		return "Vui lòng nhập email hợp lệ"
	}

	//address
	if user.ProvinceCode == "0" {
		return "Vui lòng chọn tỉnh/thành phố"
	}
	if user.DistrictCode == "0" {
		return "Vui lòng chọn quận/huyện"
	}
	if user.WardCode == "0" {
		return "Vui lòng chọn xã/phường/thị trấn"
	}

	return "ok"
}

func GetChangePasswordPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "GetChangePasswordPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	return c.Render("pages/info/change_password", fiber.Map{
		"Permissions": permissions,
		"User":        userLogin,
		"Ctx":         c,
	}, "layouts/main")
}

func PutChangePassword(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "PutChangePassword")
	var user models.User
	var form structs.FormChangePass
	DB := initializers.DB

	if err := c.BodyParser(&form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
		return c.RedirectBack("")
	}

	if form.OldPass == "" {
		return c.JSON("Mật khẩu hiện tại không được để trống")
	}
	if form.NewPass == "" {
		return c.JSON("Mật khẩu mới không được để trống")
	}
	if form.NewPass != form.CfPass {
		return c.JSON("Mật khẩu xác nhận không trùng khớp")
	}
	sess, _ := SessAuth.Get(c)
	userEmail := sess.Get("email").(string)

	if err := DB.Table("users").Where("email", userEmail).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found account")
		return c.JSON("Không tìm thấy tài khoản")
	}

	if !utils.CheckPasswordHash(form.OldPass, user.Password) {
		return c.JSON("Mật khẩu hiện tại không đúng")
	}

	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, form.NewPass)
		if !t {
			return c.JSON("Mật khẩu mới quá yếu")
		}
	}

	user.Password = utils.HashingPassword(form.NewPass)
	user.UpdatedAt = time.Now()
	user.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Table("users").Where("email", userEmail).Updates(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update password")
		return c.JSON("Đã xảy ra lỗi trong quá trình đổi mật khẩu")
	}

	return c.JSON("Success")
}

func CheckPasswordHandler(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "CheckPasswordHandler")
	// Lấy mật khẩu từ yêu cầu AJAX
	inputPassword := c.FormValue("password")

	userLogin := GetSessionUser(c)
	userPassword := userLogin.Password
	valid := utils.CheckPasswordHash(inputPassword, userPassword)

	if valid {
		return c.JSON("valid")
	}
	return c.JSON("invalid")
}
