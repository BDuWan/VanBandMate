package controllers

import (
	"github.com/zetamatta/go-outputdebug"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"vanbandmate/initializers"
	"vanbandmate/models"
	"vanbandmate/structs"
	"vanbandmate/utils"

	"github.com/gofiber/fiber/v2"
)

func GetForgotPasswordPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetForgotPasswordPage")
	return c.Render("pages/login/forgot_password", fiber.Map{})
}

func SendOtp(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: SendOTP")
	type RequestBody struct {
		Email string `json:"email"`
	}

	var body RequestBody

	// Lấy dữ liệu từ body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	userEmail := body.Email
	DB := initializers.DB
	if userEmail == "" {
		return c.JSON("Bạn chưa nhập email")
	}

	characters := "0123456789"

	var otp string

	for i := 0; i < 6; i++ {
		p := rand.Intn(len(characters))
		otp += string(characters[p])
	}

	var existingUser models.User
	if err := DB.Where("email = ?", userEmail).First(&existingUser).Error; err != nil {
		return c.JSON("Không tìm thấy tài khoản")
	}

	existingUser.Otp = otp
	existingUser.TimeExpiredOtp = time.Now().Add(10 * time.Minute)
	if err := DB.Where("email = ?", userEmail).Updates(&existingUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]" + err.Error())
		return c.JSON("Đã xảy ra lỗi trong quá trình tạo OTP")
	}

	subject := "Quên mật khẩu"
	bodyEmail := "Sử dụng mã OTP sau để cập nhật mật khẩu, mã sẽ hết hạn sau 10 phút: " + otp

	sendEmailSuccess := SendEmail(subject, bodyEmail, existingUser.Email)
	if sendEmailSuccess {
		return c.JSON("success")
	} else {
		return c.JSON("Gửi email thất bại")
	}
}

func PutUpdatePassword(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "PutUpdatePassword")
	var user models.User
	var form structs.FormUpdatePass
	DB := initializers.DB

	if err := c.BodyParser(&form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
		return c.RedirectBack("")
	}

	if strings.TrimSpace(form.Email) == "" {
		return c.JSON("Bạn chưa nhập email")
	}

	if form.Password == "" {
		return c.JSON("Mật khẩu không được để trống")
	}
	if form.Password != form.CfPassword {
		return c.JSON("Mật khẩu xác nhận không trùng khớp")
	}

	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, form.Password)
		if !t {
			return c.JSON("Mật khẩu quá yếu")
		}
	}

	if strings.TrimSpace(form.Otp) == "" {
		return c.JSON("Vui lòng nhập mã xác nhận")
	}

	if err := DB.Table("users").Where("email", form.Email).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found account")
		return c.JSON("Không tìm thấy tài khoản")
	}

	if form.Otp != user.Otp {
		return c.JSON("Mã xác nhận không đúng")
	}

	if time.Now().After(user.TimeExpiredOtp) {
		return c.JSON("Mã xác nhận đã hết hạn")
	}

	user.Password = utils.HashingPassword(form.Password)
	user.UpdatedAt = time.Now()
	user.UpdatedBy = GetSessionUser(c).UserID

	if err := DB.Table("users").Where("email", form.Email).Updates(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update password")
		return c.JSON("Đã xảy ra lỗi trong quá trình đổi mật khẩu")
	}

	return c.JSON("success")
}
