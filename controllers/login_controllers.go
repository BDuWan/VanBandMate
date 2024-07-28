package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"image"
	"image/jpeg"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"lms/utils"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/zetamatta/go-outputdebug"
)

var SessAuth = session.New(session.Config{
	CookieSessionOnly: true,
})

const (
	secretKey = "6Le0zQoqAAAAAK84axABs1IwZBcloQygggCFowuK"
)

func GetLogin(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetLogin")
	return c.Render("pages/login/login", fiber.Map{})
}
func GetSignup(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetSignup")
	DB := initializers.DB
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward
	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := DB.Find(&districts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if err := DB.Find(&wards).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	return c.Render("pages/login/signup", fiber.Map{
		"provinces": provinces,
		"districts": districts,
		"wards":     wards,
	})
}
func GetSignupSale(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetSignupSale")
	return c.Render("pages/login/signup_sale", fiber.Map{})
}
func PostLogin(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostLogin")
	var user models.User

	DB := initializers.DB
	form := new(models.User)

	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.JSON("Can not get input data")
	}

	if form.Email == "" {
		return c.JSON("Email can not be blank")
	}
	if form.Password == "" {
		return c.JSON("Password can not be blank")
	}
	if err := DB.Where(
		"BINARY email = ?", form.Email).Where(
		"deleted", false).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Email does not exist")
	}

	if !user.State {
		return c.JSON("Can not access. Please wait for the administrator to approve the account.")
	}

	if user.RoleID == 2 {
		return c.JSON("This is a sales director account. Please log in at sales department")
	}

	if !utils.CheckPasswordHash(form.Password, user.Password) {
		return c.JSON("Wrong password")
	}
	user.Session = "session_" + form.Email
	if err := DB.Save(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("")
	}

	sess, _ := SessAuth.Get(c)
	sess.Set("email", user.Email)
	sess.Set("login_success", "authenticated")
	sess.Set("user_id", user.UserID)
	sess.Set("role_id", user.RoleID)
	sess.Set("sessionId", "session_"+form.Email)
	if err := sess.Save(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON("Login success")
}

func PostSignup(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostSignup")
	var signupForm structs.SignUpForm
	var account models.User
	DB := initializers.DB
	if err := c.BodyParser(&signupForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}

	//process signup logic
	validator := ValidatorSignUpInput(signupForm)
	if validator != "ok" {
		return c.JSON(validator)
	}
	//reCAPTCHA
	fmt.Println(signupForm.Email)
	recaptchaResponse := signupForm.RecaptchaResponse
	if recaptchaResponse == "" {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: reCAPTCHA response not found")
		return c.JSON("Vui lòng xác thực reCaptcha")
	}

	client := resty.New()
	resp, err := client.R().
		SetFormData(map[string]string{
			"secret":   secretKey,
			"response": recaptchaResponse,
		}).
		Post("https://www.google.com/recaptcha/api/siteverify")

	if err != nil {
		return c.JSON("Failed to verify reCAPTCHA")
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return c.JSON("Failed to parse reCAPTCHA response")
	}

	if success, ok := result["success"].(bool); !ok || !success {
		return c.JSON("reCAPTCHA verification failed")
	}

	layout := "02/01/2006"
	date, err := time.Parse(layout, signupForm.DateOfBirth)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: error parsing date of birth")
	}
	if err := DB.Where("email", signupForm.Email).First(&models.User{}).Error; err != nil {
		account.RoleID = signupForm.RoleID
		account.FirstName = signupForm.FirstName
		account.LastName = signupForm.LastName
		account.Email = signupForm.Email
		account.PhoneNumber = signupForm.PhoneNumber
		account.ProvinceCode = signupForm.ProvinceCode
		account.DistrictCode = signupForm.DistrictCode
		account.WardCode = signupForm.WardCode
		account.AddressDetail = signupForm.AddressDetail
		account.Password = utils.HashingPassword(signupForm.Password)
		account.DateOfBirth = date
		account.Session = ""
		account.State = true
		account.Verify = false
		account.Deleted = false
		account.CreatedAt = time.Now()
		account.DeletedAt = time.Now()
		account.UpdatedAt = time.Now()
		account.DeletedBy = 0
		account.UpdatedBy = 0

		if err := DB.Create(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Can not create account")
		}

		account.CreatedBy = account.UserID

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Can not create account")
		}

		var admin []string
		if err := DB.Table("users").Where("deleted", false).Where("type_user_id = 1").Pluck("email", &admin).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found admin")
		}
		for _, item := range admin {
			go SendEmail("New registered account", "Student with Email: "+account.Email+" register for a new account.", item)
		}

		if signupForm.Image != "" {
			path := "public/assets/img/avatar/"
			saveImageResult := SaveImage(signupForm.Image, path, "123.jpg")
			if saveImageResult != "ok" {
				return c.JSON(saveImageResult)
			}
		}
		return c.JSON("Success")
	}

	return c.JSON("Email already exists")
}

func GetLogout(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetLogout")
	sess, _ := SessAuth.Get(c)

	isSale := false
	if sess.Get("role_id") == 2 {
		isSale = true
	}

	if err := sess.Reset(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	if err := sess.Save(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	if isSale {
		return c.Redirect("/login_sale")
	}
	return c.Redirect("/login")
}

func ValidatorSignUpInput(user structs.SignUpForm) string {
	if user.FirstName == "" || user.LastName == "" {
		return "Họ và tên không được để trống"
	}
	if user.Email == "" {
		return "Email không được để trống"
	}
	if user.PhoneNumber == "" {
		return "Số điện thoại không được để trống"
	}

	if user.Password == "" {
		return "Mật khẩu không được để trống"
	}
	if user.Password != user.ConfirmPassword {
		return "Mật khẩu xác nhận không trùng khớp"
	}
	if user.RoleID == 0 {
		return "Bạn chưa chọn loại tài khoản"
	}
	if user.ProvinceCode == "0" {
		return "Vui lòng chọn tỉnh/thành phố"
	}
	if user.DistrictCode == "0" {
		return "Vui lòng chọn quận/huyện"
	}
	if user.WardCode == "0" {
		return "Vui lòng chọn xã/phường/thị trấn"
	}
	regexFirstname := "^[a-zA-Z]{2,}$"
	regexFN := regexp.MustCompile(regexFirstname)
	if !regexFN.MatchString(user.FirstName) {
		return "Tên không được có số hoặc kí tự đặc biệt"
	}

	regexLastname := "^[a-zA-Z]{2,}$"
	regexLN := regexp.MustCompile(regexLastname)
	if !regexLN.MatchString(user.LastName) {
		return "Họ và tên đệm không được có số hoặc kí tự đặc biệt"
	}

	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexE := regexp.MustCompile(regexEmail)
	if !regexE.MatchString(user.Email) {
		return "Vui lòng nhập email hợp lệ"
	}

	regexPhone := "^[0-9]{10,}$"
	regexP := regexp.MustCompile(regexPhone)
	if !regexP.MatchString(user.PhoneNumber) {
		return "Vui lòng nhập số điện thoại hợp lệ"
	}

	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, user.Password)
		if !t {
			return "Mật khẩu cần có tối thiểu 8 kí tự, phải bao gồm ít nhất 1 chữ hoa, chữ thường, số và kí tự đặc biệt"
		}
	}
	return "ok"
}

func SaveImage(image_base64 string, path string, image_name string) string {
	os.Remove(path + image_name)

	imageBase64 := image_base64

	data := strings.Split(imageBase64, ",")[1]

	imageBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Cannot upload image"
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Cannot upload image"
	}

	out, err := os.Create(path + image_name)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Cannot upload image"
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Cannot upload image"
	}
	return "ok"
}
