package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gopkg.in/gomail.v2"
	"os"
	"time"
	"vanbandmate/initializers"
	"vanbandmate/models"
)

var key = []byte(os.Getenv("KEY_TOKEN"))

func createToken(email string) (string, error) {

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateToken(tokenString string) (models.User, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return models.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		currentTime := time.Now()

		if currentTime.After(expirationTime) {
			return models.User{}, fmt.Errorf("token has expired")
		}

		user := models.User{
			Email: claims["email"].(string),
		}
		return user, nil
	}

	return models.User{}, fmt.Errorf("invalid token")
}

func SendEmail(subject, body string, emailUser string) bool {
	email := "vanbandmate@gmail.com"
	password := "tdpm ozap iuwb wpvy"

	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	to := emailUser

	message := gomail.NewMessage()
	message.SetHeader("From", email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "VanBandMate - "+subject)
	message.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, email, password)

	if err := d.DialAndSend(message); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]:" + err.Error())
		return false
	}
	return true
}

func VerifyHandler(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: VerifyHander")
	token := c.Query("token")

	user, err := validateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}
	DB := initializers.DB

	var existingUser models.User
	if err := DB.Table("users").Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return c.JSON("User not found")
	}

	existingUser.Verify = true

	DB.Table("users").Where("email = ?", user.Email).Updates(&existingUser)

	return c.Redirect("/info")
}

func PostVerifyEmail(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostVerifyEmail")
	sess, _ := SessAuth.Get(c)
	DB := initializers.DB
	var user models.User
	userEmail, ok := sess.Get("email").(string)
	if !ok {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Cannot get user from session")
		return c.RedirectBack("")
	}

	if err := DB.Where("email", userEmail).First(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Not found account")
		return c.JSON("Đã xảy ra lỗi khi tìm kiếm tài khoản")
	}
	token, err := createToken(user.Email)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		// Xử lý lỗi nếu cần
		return c.JSON("Đã xảy ra lỗi khi tạo token")
	}

	subject := "Xác thực email"
	bodyEmail := "Ấn vào link để xác thực, link sẽ hết hạn sau 24h " + fmt.Sprintf(os.Getenv("HOST")+"/verify?token=%s", token)
	sendEmailSuccess := SendEmail(subject, bodyEmail, user.Email)
	if sendEmailSuccess {
		return c.JSON("success")
	} else {
		return c.JSON("Gửi email thất bại")
	}
}
