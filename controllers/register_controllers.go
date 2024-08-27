package controllers

import (
	"fmt"
	"lms/models"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func IsChecked(currentID int, checkedID []models.RolePermission) bool {

	for _, check := range checkedID {
		if check.PermissionID == currentID {
			return true
		}
	}

	return false
}

func IsSelected(currentID, selectedID int) bool {
	return currentID == selectedID
}

func IsVerify(c *fiber.Ctx) bool {

	return GetSessionUser(c).Verify
}

//	GetUserID func IsStudent(typeUserID int) bool {
//		if typeUserID == 5 {
//			return true
//		}
//
//		return false
//	}
func GetUserID(c *fiber.Ctx) int {
	sess, _ := SessAuth.Get(c)
	idUser := fmt.Sprintf("%d", sess.Get("user_id"))
	userId, _ := strconv.Atoi(idUser)

	return userId
}
func FormatDate(t time.Time) string {
	return t.Format("02/01/2006")
}
func FormatTime(eventTime time.Time) string {
	return eventTime.Format("2006-01-02 15:04:05")
}
func IsTimeAfterNow(t time.Time) bool {
	return t.After(time.Now())
}
func FormatFloat64(f float64) float64 {
	scale := math.Pow10(3)
	return math.Round(f*scale) / scale
}
func FormatTimeComment(t time.Time) string {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	if t.After(today) {
		return t.Format("15:04")
	}
	return t.Format("01-02-2006")
}
