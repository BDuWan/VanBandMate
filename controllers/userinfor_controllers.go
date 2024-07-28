package controllers

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
//func GetInformation(c *fiber.Ctx) error {
//	var user models.User
//	var roles []models.Role
//	var typeUsers []models.TypeUser
//
//	DB := initializers.DB
//	sess, _ := SessAuth.Get(c)
//	sessionID, ok := sess.Get("sessionId").(string)
//	if !ok {
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).First(&user).Error; err != nil {
//		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//			"message": "User not found",
//		})
//	}
//
//	if err := DB.Where("deleted", false).Find(&roles).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		fmt.Println("loi3")
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Select("type_user_id", "name").Find(&typeUsers).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		fmt.Println("loi4")
//		return c.RedirectBack("")
//	}
//
//	//data := struct {
//	//	Roles      []models.Role
//	//	UserRoleID int
//	//}{
//	//	Roles:      roles,
//	//	UserRoleID: role.RoleID, // Set the ID of the selected category
//	//}
//
//	return c.Render("pages/home/information", fiber.Map{
//		"User": user,
//		"Ctx":  c,
//	}, "layouts/main")
//}
//
//func PutUpdateUserInformation(c *fiber.Ctx) error {
//	var user models.User
//	var form structs.AccUser
//	DB := initializers.DB
//
//	if err := c.BodyParser(&form); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User Fail")
//	}
//	sess, _ := SessAuth.Get(c)
//	sessionID, ok := sess.Get("sessionId").(string)
//	if !ok {
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).First(&user).Error; err != nil {
//		return c.JSON("Can not found account")
//	}
//	if user.Verify == false {
//		return c.JSON("Unable to update information, please verify account")
//	}
//
//	user.FirstName = form.FirstName
//	user.LastName = form.LastName
//	user.RoleID = form.RoleID
//	user.Email = form.Email
//	user.PhoneNumber = form.PhoneNumber
//	user.Address = form.Address
//	user.UpdatedAt = time.Now()
//	user.UpdatedBy = GetSessionUser(c).UserID
//
//	var checkRef models.User
//	// referralCode := true
//	if user.ReferralCode == "" {
//		if form.ReferralCode != "" {
//			if err := DB.Where("deleted", false).Where("code_user", user.ReferralCode).First(&checkRef).Error; err != nil {
//				// referralCode = false
//				user.ReferralCode = form.ReferralCode
//			}
//		}
//
//	}
//
//	result := ValidatorUpdateAccUser(form)
//	if result != "ok" {
//		return c.JSON(result)
//	}
//
//	if err := DB.Table("users").Where("session = ?", sessionID).Updates(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update account")
//		return c.JSON("Can not update account")
//	}
//
//	// if user.TypeUserID == 5 {
//	// 	if referralCode == false {
//	// 		resultGainCommission := GainCommissionSaleBusiness(user.ReferralCode, user.CreatedAt)
//	// 		if resultGainCommission != "ok" {
//	// 			return c.SendString(resultGainCommission)
//	// 		}
//	// 	}
//
//	// }
//	return c.JSON("Success")
//}
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
