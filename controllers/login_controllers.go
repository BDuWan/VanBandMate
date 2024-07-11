package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"lms/utils"
	"regexp"
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
func GetLoginSale(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetLoginSale")
	return c.Render("pages/login/login_sale", fiber.Map{})
}
func GetSignup(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetSignup")
	return c.Render("pages/login/signup", fiber.Map{})
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

func PostLoginSale(c *fiber.Ctx) error {
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

	if user.RoleID != 2 {
		return c.JSON("Can not access. This is not a sales director account")
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

func PostSignupStudent(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostSignupStudent")
	var user structs.User
	var account models.User
	DB := initializers.DB
	if err := c.BodyParser(&user); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}
	//reCAPTCHA
	recaptchaResponse := user.RecaptchaResponse
	if recaptchaResponse == "" {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: reCAPTCHA response not found")
		return c.JSON("reCAPTCHA response not found")
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

	//process signup logic
	validator := ValidatorSignUpInput(user)
	if validator != "ok" {
		return c.JSON(validator)
	}

	if err := DB.Where("email", user.Email).First(&models.User{}).Error; err != nil {
		var checkReferralCode models.User
		if user.ReferralCode != "" {
			if err := DB.Where("type_user_id = ? or type_user_id = ?", 2, 3).Where(
				"code_user", user.ReferralCode).First(&checkReferralCode).Error; err != nil {
				return c.JSON("Please enter the correct referral code or leave it blank")
			}
		}

		account.TypeUserID = 5
		account.FirstName = user.FirstName
		account.LastName = user.LastName
		account.RoleID = 5
		account.Email = user.Email
		account.PhoneNumber = user.PhoneNumber
		account.Address = user.Address
		account.Username = user.Email
		account.Password = utils.HashingPassword(user.Password)
		account.ReferralCode = user.ReferralCode
		account.Session = ""
		account.NameBusiness = ""
		account.FullNameRepresentative = ""
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

		account.CodeUser = "STUD" + fmt.Sprintf("%04d", account.UserID)
		account.CreatedBy = account.UserID

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Can not create account")
		}
		if user.ReferralCode != "" {
			resultGainNumberStudentSaleBusiness := GainNumberStudentSaleBusiness(account.UserID, user.ReferralCode, account.CreatedAt)
			if resultGainNumberStudentSaleBusiness != "ok" {
				return c.SendString(resultGainNumberStudentSaleBusiness)
			}
		}

		var admin []string
		if err := DB.Table("users").Where("deleted", false).Where("type_user_id = 1").Pluck("email", &admin).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found admin")
		}
		for _, item := range admin {
			go SendEmail("New registered account", "Student with Email: "+account.Email+" register for a new account.", item)
		}
		return c.JSON("Success")
	}

	return c.JSON("Email already exists")
}

func PostSignupInstructor(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostSignupInstructor")
	var user structs.User
	var account models.User
	DB := initializers.DB
	if err := c.BodyParser(&user); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}
	//reCAPTCHA
	recaptchaResponse := user.RecaptchaResponse
	if recaptchaResponse == "" {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: reCAPTCHA response not found")
		return c.JSON("reCAPTCHA response not found")
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

	//process signup logic
	validator := ValidatorSignUpInput(user)
	if validator != "ok" {
		return c.JSON(validator)
	}

	if err := DB.Where("email", user.Email).First(&models.User{}).Error; err != nil {
		account.TypeUserID = 4
		account.FirstName = user.FirstName
		account.LastName = user.LastName
		account.RoleID = 4
		account.Email = user.Email
		account.PhoneNumber = user.PhoneNumber
		account.Address = user.Address
		account.Username = user.Email
		account.Password = utils.HashingPassword(user.Password)
		account.ReferralCode = ""
		account.Session = ""
		account.NameBusiness = ""
		account.FullNameRepresentative = ""
		account.State = false
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

		account.CodeUser = "INSTR" + fmt.Sprintf("%04d", account.UserID)
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
			go SendEmail("New registered account", "Instructor with Email: "+account.Email+" register for a new account. Please approve this account.", item)
		}
		return c.JSON("Success")
	}

	return c.JSON("Email already exists")
}

func PostSignupSale(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostSignupSale")
	var user structs.User
	var account models.User
	DB := initializers.DB
	if err := c.BodyParser(&user); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}
	//reCAPTCHA
	recaptchaResponse := user.RecaptchaResponse
	if recaptchaResponse == "" {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: reCAPTCHA response not found")
		return c.JSON("reCAPTCHA response not found")
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

	//process signup logic

	validator := ValidatorSignUpInput(user)
	if validator != "ok" {
		return c.JSON(validator)
	}

	if err := DB.Where("email", user.Email).First(&models.User{}).Error; err != nil {
		account.TypeUserID = 2
		account.FirstName = user.FirstName
		account.LastName = user.LastName
		account.RoleID = 2
		account.Email = user.Email
		account.PhoneNumber = user.PhoneNumber
		account.Address = user.Address
		account.Username = user.Email
		account.Password = utils.HashingPassword(user.Password)
		account.ReferralCode = ""
		account.Session = ""
		account.NameBusiness = ""
		account.FullNameRepresentative = ""
		account.State = false
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

		account.CodeUser = "SAL" + fmt.Sprintf("%04d", account.UserID)
		account.CreatedBy = account.UserID

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Can not create account")
		}
		var commissionUser models.CommissionUser
		commissionUser.UserID = account.UserID
		commissionUser.CommissionTotal = 0
		commissionUser.CommissionPaid = 0
		commissionUser.CommissionDebt = 0
		commissionUser.PeriodID = 0
		commissionUser.Deleted = false
		commissionUser.CreatedAt = time.Now()
		commissionUser.CreatedBy = account.UserID
		if err := DB.Create(&commissionUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create commission user")
		}
		var admin []string
		if err := DB.Table("users").Where("deleted", false).Where("type_user_id = 1").Pluck("email", &admin).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found admin")
		}
		for _, item := range admin {
			go SendEmail("New registered account", "Sale business with Email: "+account.Email+" register for a new account. Please approve this account.", item)
		}
		return c.JSON("Success")
	}

	return c.JSON("Email already exists")
}

func PostSignupBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostSignupBusiness")
	var user structs.User
	var account models.User
	DB := initializers.DB
	if err := c.BodyParser(&user); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
	}

	validator := ValidatorSignUpBusinessInput(user)
	if validator != "ok" {
		return c.JSON(validator)
	}

	if err := DB.Where("email", user.Email).First(&models.User{}).Error; err != nil {
		account.TypeUserID = 3
		account.FirstName = ""
		account.LastName = ""
		account.RoleID = 3
		account.Email = user.Email
		account.PhoneNumber = user.PhoneNumber
		account.Address = user.Address
		account.Username = user.Email
		account.Password = utils.HashingPassword(user.Password)
		account.ReferralCode = ""
		account.Session = ""
		account.NameBusiness = user.NameBusiness
		account.FullNameRepresentative = user.FullNameRepresentative
		account.State = false
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

		account.CodeUser = "BUSIN" + fmt.Sprintf("%04d", account.UserID)
		account.CreatedBy = account.UserID

		if err := DB.Updates(&account).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
			return c.JSON("Can not create account")
		}
		var commissionUser models.CommissionUser
		commissionUser.UserID = account.UserID
		commissionUser.CommissionTotal = 0
		commissionUser.CommissionPaid = 0
		commissionUser.CommissionDebt = 0
		commissionUser.PeriodID = 0
		commissionUser.Deleted = false
		commissionUser.CreatedAt = time.Now()
		commissionUser.CreatedBy = account.UserID
		if err := DB.Create(&commissionUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create commission user")
		}
		var admin []string
		if err := DB.Table("users").Where("deleted", false).Where("type_user_id = 1").Pluck("email", &admin).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Not found admin")
		}
		for _, item := range admin {
			go SendEmail("New registered account", "Sale business with Email: "+account.Email+" register for a new account. Please approve this account.", item)
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

func ValidatorSignUpInput(user structs.User) string {
	if user.FirstName == "" {
		return "Firstname can not be blank"
	}
	if user.LastName == "" {
		return "Lastname can not be blank"
	}
	if user.Email == "" {
		return "Email can not be blank"
	}
	if user.PhoneNumber == "" {
		return "PhoneNumber can not be blank"
	}
	if user.Address == "" {
		return "Address can not be blank"
	}
	//if user.Username == "" {
	//	return "Username can not be blank"
	//}
	if user.Password == "" {
		return "Password can not be blank"
	}
	//regexFirstname := "^[a-zA-Z]{2,}$"
	//regexFN := regexp.MustCompile(regexFirstname)
	//if !regexFN.MatchString(user.FirstName) {
	//	return "Invalid First name"
	//}
	//
	//regexLastname := "^[a-zA-Z]{2,}$"
	//regexLN := regexp.MustCompile(regexLastname)
	//if !regexLN.MatchString(user.LastName) {
	//	return "Invalid Last name"
	//}

	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexE := regexp.MustCompile(regexEmail)
	if !regexE.MatchString(user.Email) {
		return "Invalid Email"
	}

	regexPhone := "^[0-9]{7,}$"
	regexP := regexp.MustCompile(regexPhone)
	if !regexP.MatchString(user.PhoneNumber) {
		return "Invalid Phone number"
	}

	//regexUsername := "^[a-zA-Z0-9]{8,}$"
	//regexU := regexp.MustCompile(regexUsername)
	//if !regexU.MatchString(user.Username) {
	//	return "Username must have at least 8 characters, including only letters and numbers"
	//}

	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, user.Password)
		if !t {
			return "Password too weak"
		}
	}
	return "ok"
}

func ValidatorSignUpBusinessInput(user structs.User) string {
	if user.NameBusiness == "" {
		return "Name of business can not be blank"
	}
	if user.FullNameRepresentative == "" {
		return "Full name of representative can not be blank"
	}
	if user.Email == "" {
		return "Business email can not be blank"
	}
	if user.PhoneNumber == "" {
		return "Business phoneNumber can not be blank"
	}
	if user.Address == "" {
		return "Business address can not be blank"
	}
	//if user.Username == "" {
	//	return "Business username can not be blank"
	//}
	if user.Password == "" {
		return "Business password can not be blank"
	}
	regexNameBusiness := "^[a-zA-Z0-9]{2,}$"
	regexNB := regexp.MustCompile(regexNameBusiness)
	if !regexNB.MatchString(user.NameBusiness) {
		return "Invalid Name Business"
	}

	regexFullNameR := "^[a-zA-Z0-9]{2,}$"
	regexFNR := regexp.MustCompile(regexFullNameR)
	if !regexFNR.MatchString(user.FullNameRepresentative) {
		return "Invalid Full Name Representative"
	}

	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexE := regexp.MustCompile(regexEmail)
	if !regexE.MatchString(user.Email) {
		return "Invalid Email"
	}

	regexPhone := "^[0-9]{7,}$"
	regexP := regexp.MustCompile(regexPhone)
	if !regexP.MatchString(user.PhoneNumber) {
		return "Invalid Phone number"
	}

	//regexUsername := "^[a-zA-Z0-9]{8,}$"
	//regexU := regexp.MustCompile(regexUsername)
	//if !regexU.MatchString(user.Username) {
	//	return "Business username must have at least 8 characters, including only letters and numbers"
	//}

	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
	for _, test := range testPassword {
		t, _ := regexp.MatchString(test, user.Password)
		if !t {
			return "Business password too weak"
		}
	}
	return "ok"
}
