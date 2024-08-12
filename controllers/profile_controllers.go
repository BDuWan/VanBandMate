package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//
//import (
//	"fmt"
//	"lms/initializers"
//	"lms/models"
//	"lms/structs"
//	"lms/utils"
//	"os"
//	"regexp"
//	"time"
//
//	"github.com/dgrijalva/jwt-go"
//	"github.com/gofiber/fiber/v2"
//	"github.com/zetamatta/go-outputdebug"
//	"gopkg.in/gomail.v2"
//)
//

func GetProfileID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Get Profile ID")
	userId := c.Params("id")
	DB := initializers.DB
	user := new(models.User)
	if err := DB.Model(&models.User{}).Joins("Role").Joins(
		"Province").Joins("District").Joins("Ward").Where(
		"user_id", userId).First(user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Cannot get user")
	}
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	return c.Render("pages/home/profile", fiber.Map{
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
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Cannot get user")
	}
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")

	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := DB.Find(&districts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := DB.Find(&wards).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.Redirect("/errors/404")
	}

	return c.Render("pages/home/userinfo", fiber.Map{
		"Permissions": permissions,
		"User":        user,
		"provinces":   provinces,
		"districts":   districts,
		"wards":       wards,
		"Ctx":         c,
	}, "layouts/main")
}

func PutUserInfo(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Put User Info")
	var updateInfoForm structs.UpdateInfoForm
	var account models.User
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	if err := c.BodyParser(&updateInfoForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}
	validator := ValidatorUpdateInfoInput(updateInfoForm)
	if validator != "ok" {
		return c.JSON(validator)
	}
	layout := "02/01/2006"
	date, err := time.Parse(layout, updateInfoForm.DateOfBirth)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: error parsing date of birth")
	}

	if err := DB.Where("user_id", userLogin.UserID).First(&account).Error; err == nil {
		if account.Email != updateInfoForm.Email {
			if err := DB.Where("email", updateInfoForm.Email).First(&models.User{}).Error; err == nil {
				return c.JSON("Email đã được sử dụng, vui lòng chọn 1 email khác")
			}
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
		account.State = false
		account.Verify = false
		account.UpdatedAt = time.Now()
		account.UpdatedBy = 0

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Không thể cập nhật thông tin")
		}
		imageName := "avatar" + strconv.Itoa(account.UserID) + ".jpg"

		if updateInfoForm.Image != "" {
			path := "public/assets/img/avatar/"
			saveImageResult := SaveImage(updateInfoForm.Image, path, imageName)
			if saveImageResult != "ok" {
				return c.JSON(saveImageResult)
			}
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

//
//func ValidatorUpdateAccUser(user structs.AccUser) string {
//	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
//	regexE := regexp.MustCompile(regexEmail)
//	if !regexE.MatchString(user.Email) {
//		return "Invalid Email"
//	}
//
//	regexFirstname := "^[a-zA-Z]{2,}$"
//	regexFN := regexp.MustCompile(regexFirstname)
//	if !regexFN.MatchString(user.FirstName) {
//		return "Invalid First name"
//	}
//
//	regexLastname := "^[a-zA-Z]{2,}$"
//	regexLN := regexp.MustCompile(regexLastname)
//	if !regexLN.MatchString(user.LastName) {
//		return "Invalid Last name"
//	}
//
//	regexPhone := "^[0-9]{7,}$"
//	regexP := regexp.MustCompile(regexPhone)
//	if !regexP.MatchString(user.PhoneNumber) {
//		return "Invalid Phone number"
//	}
//	return "ok"
//}
//func GetChangePassword(c *fiber.Ctx) error {
//	return c.Render("pages/home/changePassword", fiber.Map{
//		"Ctx": c,
//	}, "layouts/main")
//}
//
//func PutChangePassword(c *fiber.Ctx) error {
//	var user models.User
//	var form structs.FromChangePass
//	DB := initializers.DB
//
//	if err := c.BodyParser(&form); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
//		return c.RedirectBack("")
//	}
//
//	if form.CurrentPass == "" {
//		return c.JSON("Current password can not be blank")
//	}
//	if form.NewPass == "" {
//		return c.JSON("New password can not be blank")
//	}
//	sess, _ := SessAuth.Get(c)
//	sessionID, ok := sess.Get("sessionId").(string)
//	if !ok {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Cannot get user from session")
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).First(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found account")
//		return c.JSON("Not found account")
//	}
//
//	if !utils.CheckPasswordHash(form.CurrentPass, user.Password) {
//		return c.JSON("Current password not match")
//	}
//
//	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
//	for _, test := range testPassword {
//		t, _ := regexp.MatchString(test, form.NewPass)
//		if !t {
//			return c.JSON("New password too weak")
//		}
//	}
//
//	user.Password = utils.HashingPassword(form.NewPass)
//	user.UpdatedAt = time.Now()
//	user.UpdatedBy = GetSessionUser(c).UserID
//
//	if err := DB.Table("users").Where("session = ?", sessionID).Updates(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update password")
//		return c.JSON("Can not update password")
//	}
//
//	return c.JSON("Success")
//}
//
//var key = []byte("1234")
//
//func createAndSaveToken(email string) error {
//
//	claims := jwt.MapClaims{
//		"email": email,
//		"exp":   time.Now().Add(time.Hour * 24).Unix(),
//	}
//	DB := initializers.DB
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	signedToken, err := token.SignedString(key)
//	if err != nil {
//		return err
//	}
//
//	var updateUser models.User
//
//	updateUser.Token = signedToken
//
//	if err := DB.Table("users").Where("email = ?", email).Updates(&updateUser); err != nil {
//		return fmt.Errorf("can not save token")
//	}
//
//	return nil
//}
//
//func validateToken(tokenString string) (models.User, error) {
//
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
//		}
//		return key, nil
//	})
//
//	if err != nil {
//		return models.User{}, err
//	}
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//
//		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
//		currentTime := time.Now()
//
//		if currentTime.After(expirationTime) {
//			return models.User{}, fmt.Errorf("token has expired")
//		}
//
//		user := models.User{
//			Email: claims["email"].(string),
//		}
//		return user, nil
//	}
//
//	return models.User{}, fmt.Errorf("invalid token")
//}
//
//func SendVerifyEmail(c *fiber.Ctx) string {
//	sess, _ := SessAuth.Get(c)
//	DB := initializers.DB
//	var user models.User
//	sessionID, ok := sess.Get("sessionId").(string)
//	if !ok {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Cannot get user from session")
//		return "error"
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).First(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found account")
//		return "Not found account"
//	}
//	token := user.Token
//
//	//email := "buiv03061@gmail.com"
//	//password := "usno timq hhks vbyb"
//	email := "spusefulknowledge@gmail.com"
//	password := "srjq wmyn oztk fcye"
//	smtpHost := "smtp.gmail.com"
//	smtpPort := 587
//
//	to := user.Email
//
//	subject := "Useful knowledge -  Email authentication"
//	body := "Click on the link to verify your account " + fmt.Sprintf(os.Getenv("HOST")+"/verify?token=%s", token)
//
//	message := gomail.NewMessage()
//	message.SetHeader("From", email)
//	message.SetHeader("To", to)
//	message.SetHeader("Subject", subject)
//	message.SetBody("text/plain", body)
//
//	d := gomail.NewDialer(smtpHost, smtpPort, email, password)
//
//	if err := d.DialAndSend(message); err != nil {
//		return "Email sent failed!"
//	}
//	return "Success"
//}
//
//func VerifyHandler(c *fiber.Ctx) error {
//
//	token := c.Query("token")
//
//	user, err := validateToken(token)
//	if err != nil {
//		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
//	}
//	DB := initializers.DB
//
//	var existingUser models.User
//	if err := DB.Table("users").Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
//		return c.JSON("User not found")
//	}
//
//	existingUser.Verify = true
//
//	DB.Table("users").Where("email = ?", user.Email).Updates(&existingUser)
//
//	// if existingUser.ReferralCode != "" {
//	// 	if existingUser.TypeUserID == 5 {
//	// 		resultGainCommission := GainCommissionSaleBusiness(existingUser.ReferralCode, existingUser.CreatedAt)
//	// 		if resultGainCommission != "ok" {
//	// 			return c.Redirect("/home")
//	// 		}
//	// 	}
//	// }
//
//	return c.Redirect("/home")
//	//return c.SendString(fmt.Sprintf("Email of %s has been verified", user.Username))
//}
//
//func PostVerifyEmail(c *fiber.Ctx) error {
//	sess, _ := SessAuth.Get(c)
//	DB := initializers.DB
//	var user models.User
//	sessionID, ok := sess.Get("sessionId").(string)
//	if !ok {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Cannot get user from session")
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).First(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found account")
//		return c.JSON("Not found account")
//	}
//	token := user.Token
//	if token == "" {
//		if err := createAndSaveToken(user.Email); err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		}
//	}
//	result := SendVerifyEmail(c)
//
//	return c.JSON(result)
//}
