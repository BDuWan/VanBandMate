package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"image"
	"image/jpeg"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"lms/utils"
	"os"
	"regexp"
	"strconv"
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
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetLogin")
	return c.Render("pages/login/login", fiber.Map{})
}
func GetSignup(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetSignup")
	DB := initializers.DB
	var provinces []models.Province
	var districts []models.District
	var wards []models.Ward
	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	if err := DB.Find(&districts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	if err := DB.Find(&wards).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}
	return c.Render("pages/login/signup", fiber.Map{
		"provinces": provinces,
		"districts": districts,
		"wards":     wards,
	})
}
func PostLogin(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostLogin")
	var user models.User
	//var permissions []models.RolePermission

	DB := initializers.DB
	form := new(models.User)

	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Can not get input data")
	}

	if form.Email == "" {
		return c.JSON("Vui lòng nhập email")
	}
	if form.Password == "" {
		return c.JSON("Vui lòng nhập mật khẩu")
	}
	if err := DB.Where(
		"BINARY email = ?", form.Email).Where(
		"deleted", false).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Email không tồn tại")
	}

	if !user.State {
		return c.JSON("Tài khoản chưa được admin chấp thuận. Vui lòng đợi")
	}

	if !utils.CheckPasswordHash(form.Password, user.Password) {
		return c.JSON("Sai mật khẩu")
	}
	//user.Session = "session_" + form.Email
	if err := DB.Save(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("")
	}

	//if err := DB.Model(&models.RolePermission{}).Joins("Permission").Where(
	//	"role_permissions.role_id", user.RoleID).Where(
	//	"role_permissions.deleted", false).Find(&permissions).Error; err != nil {
	//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	//}
	//fmt.Println(permissions)

	sess, _ := SessAuth.Get(c)
	sess.Set("email", user.Email)
	sess.Set("login_success", "authenticated")
	sess.Set("user_id", user.UserID)
	sess.Set("role_id", user.RoleID)
	//sess.Set("permissions", permissions)
	//sess.Set("sessionId", "session_"+form.Email)
	if err := sess.Save(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	return c.JSON("Login success")
}

func PostSignup(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostSignup")
	var signupForm structs.SignUpForm
	var account models.User
	var role models.Role
	DB := initializers.DB
	if err := c.BodyParser(&signupForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Format User Fail")
	}
	//process signup logic
	validator := ValidatorSignUpInput(signupForm)
	if validator != "ok" {
		return c.JSON(validator)
	}
	//reCAPTCHA
	recaptchaResponse := signupForm.RecaptchaResponse
	if recaptchaResponse == "" {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: reCAPTCHA response not found")
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
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: error parsing date of birth")
	}
	if err := DB.Where("email", signupForm.Email).First(&models.User{}).Error; err != nil {
		account.RoleID = signupForm.RoleID
		account.Gender = signupForm.Gender
		account.FirstName = signupForm.FirstName
		account.LastName = signupForm.LastName
		account.Email = signupForm.Email
		account.PhoneNumber = signupForm.PhoneNumber
		account.LinkFacebook = signupForm.LinkFacebook
		account.ProvinceCode = signupForm.ProvinceCode
		account.DistrictCode = signupForm.DistrictCode
		account.WardCode = signupForm.WardCode
		account.AddressDetail = signupForm.AddressDetail
		account.Password = utils.HashingPassword(signupForm.Password)
		account.DateOfBirth = date
		//account.Session = ""
		account.State = false
		account.Verify = false
		account.Deleted = false
		account.CreatedAt = time.Now()
		account.DeletedAt = time.Now()
		account.UpdatedAt = time.Now()
		account.DeletedBy = 0
		account.UpdatedBy = 0

		if err := DB.Create(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Can not create account")
			return c.JSON("Can not create account")
		}

		account.CreatedBy = account.UserID
		imageName := "avatar" + strconv.Itoa(account.UserID) + ".jpg"
		account.Image = imageName

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Can not create account")
			return c.JSON("Đã xa ra lỗi khi tạo tài khoản")
		}
		if signupForm.Image != "" {
			path := "public/assets/img/avatar/"
			saveImageResult := SaveImage(signupForm.Image, path, imageName)
			if saveImageResult != "ok" {
				return c.JSON(saveImageResult)
			}
		}

		//var admin []string
		//if err := DB.Table("users").Where("deleted", false).Where("type_user_id = 1").Pluck("email", &admin).Error; err != nil {
		//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Not found admin")
		//}
		//for _, item := range admin {
		//	go SendEmail("New registered account", "Student with Email: "+account.Email+" register for a new account.", item)
		//}
		if err := DB.Where("role_id", account.RoleID).Where("deleted", false).First(&role).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Không tìm thấy vai trò")
		}
		role.NumberUser += 1
		if err := DB.Save(&role).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Lỗi khi cập nhật vai trò")
		}
		return c.JSON("Success")
	}

	return c.JSON("Email already exists")
}

func GetLogout(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetLogout")
	sess, _ := SessAuth.Get(c)

	if err := sess.Reset(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	if err := sess.Save(); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	return c.Redirect("/login")
}

func ValidatorSignUpInput(user structs.SignUpForm) string {
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

	//role
	if user.RoleID == 0 {
		return "Bạn chưa chọn loại tài khoản"
	}

	//email
	if strings.TrimSpace(user.Email) == "" {
		return "Email không được để trống"
	}
	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexE := regexp.MustCompile(regexEmail)
	if !regexE.MatchString(user.Email) {
		return "Vui lòng nhập email hợp lệ"
	}

	//password
	if strings.TrimSpace(user.Password) == "" {
		return "Mật khẩu không được để trống"
	}
	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, user.Password)
		if !t {
			return "Mật khẩu cần có tối thiểu 8 kí tự, phải bao gồm ít nhất 1 chữ hoa, chữ thường, số và kí tự đặc biệt"
		}
	}
	if user.Password != user.ConfirmPassword {
		return "Mật khẩu xác nhận không trùng khớp"
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

func SaveImage(image_base64 string, path string, image_name string) string {
	os.Remove(path + image_name)

	imageBase64 := image_base64

	data := strings.Split(imageBase64, ",")[1]

	imageBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return "Cannot upload image"
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return "Cannot upload image"
	}

	out, err := os.Create(path + image_name)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return "Cannot upload image"
	}
	defer out.Close()

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return "Cannot upload image"
	}
	return "ok"
}
