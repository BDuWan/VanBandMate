package controllers

import (
	"lms/initializers"
	"lms/models"
	"lms/utils"
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

func CreateNewPass(c *fiber.Ctx) error {
	DB := initializers.DB
	form := new(models.User)
	err := c.BodyParser(form)
	if err != nil {
		return c.JSON("Format data failed")
	}

	if form.Email == "" {
		return c.JSON("Email blank!")
	}

	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var newPass string

	for i := 0; i < 8; i++ {
		p := rand.Intn(len(characters))
		newPass += string(characters[p])
	}

	hashedPassword := utils.HashingPassword(newPass)

	var existingUser models.User
	if err := DB.Table("users").Where("email = ?", form.Email).First(&existingUser).Error; err != nil {
		return c.JSON("Not found email")
	}

	existingUser.Password = hashedPassword
	DB.Table("users").Where("email = ?", form.Email).Updates(&existingUser)

	if SendEmail("New password", "Your new password is "+newPass, existingUser.Email) {
		return c.JSON("Send email success")
	} else {
		return c.JSON("Send email fail")
	}
}

func SendEmail(subject, body string, emailUser string) bool {
	//email := "buiv03061@gmail.com"
	//password := "usno timq hhks vbyb"
	email := "spusefulknowledge@gmail.com"
	password := "srjq wmyn oztk fcye"

	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	to := emailUser

	message := gomail.NewMessage()
	message.SetHeader("From", email)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Useful knowledge - "+subject)
	message.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, email, password)

	if err := d.DialAndSend(message); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
