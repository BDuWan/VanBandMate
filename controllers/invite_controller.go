package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"time"
)

// GetHiringInvitePage func GetFindMusicPlayerPage(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetFindMusicPlayerPage")
//	userLogin := GetSessionUser(c)
//	sess, _ := SessAuth.Get(c)
//	permissions := sess.Get("rolePermission")
//	var allPermissions []models.Permission
//	DB := initializers.DB
//	var provinces []models.Province
//	var districts []models.District
//	var wards []models.Ward
//
//	if err := DB.Find(&allPermissions).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//	}
//
//	if err := DB.Find(&provinces).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//		return c.Redirect("/errors/404")
//	}
//	if err := DB.Find(&districts).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//		return c.Redirect("/errors/404")
//	}
//	if err := DB.Find(&wards).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
//		return c.Redirect("/errors/404")
//	}
//	return c.Render("pages/find-music-player/index", fiber.Map{
//		"This":           8,
//		"Permissions":    permissions,
//		"AllPermissions": allPermissions,
//		"User":           userLogin,
//		"provinces":      provinces,
//		"districts":      districts,
//		"wards":          wards,
//		"Ctx":            c,
//	}, "layouts/main")
//}

func GetHiringInvitePage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetHiringInvitePage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB

	var hiringNewsId = c.Params("id")
	var hiringNews models.HiringNews

	if err := DB.Model(&models.HiringNews{}).
		Joins("Province").Joins("District").Joins("Ward").
		Where("hiring_news_id", hiringNewsId).Where("deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	if hiringNews.ChuloadaiID != userLogin.UserID {
		return c.Redirect("/errors/403")
	}

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	return c.Render("pages/hiring/invitation", fiber.Map{
		"This":           7,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"User":           userLogin,
		"HiringNews":     hiringNews,
		"Ctx":            c,
	}, "layouts/main")
}

func APIPostHiringFind(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostHiringFind")
	DB := initializers.DB

	form := new(structs.FormFind)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	hiringNews := form.HiringNews
	var nhaccongIds []int

	if err := DB.Model(&models.Contract{}).
		Where("DATE(date) = ?", hiringNews.Date.Format("2006-01-02")).
		Where("deleted = ?", false).
		Pluck("nhaccong_id", &nhaccongIds).
		Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	page := form.Page
	limit := form.ItemsPerPage

	offset := (page - 1) * limit

	userLogin := GetSessionUser(c)
	query := DB.Model(&models.User{}).
		Joins("Province").Joins("District").Joins("Ward").
		Where("users.deleted", false).
		Where("users.get_invitation", true)

	if len(nhaccongIds) > 0 {
		query = query.Where("users.user_id NOT IN ?", nhaccongIds)
	}
	if form.Condition == 1 {
		query = query.Where("users.province_code", hiringNews.ProvinceCode)
	}
	if form.Condition == 2 {
		query = query.Where("users.district_code", hiringNews.DistrictCode)
	}
	if form.Condition == 3 {
		query = query.Where("users.ward_code", hiringNews.WardCode)
	}

	var totalItems int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&totalItems).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var nhaccongs []models.User
	if err := query.Offset(offset).Limit(limit).Find(&nhaccongs).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	var invitations []models.Invitation
	if err := DB.
		Where("hiring_news_id", hiringNews.HiringNewsID).
		Find(&invitations).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	invitationMap := make(map[int]int)
	for _, inv := range invitations {
		invitationMap[inv.NhaccongID] = inv.Status
	}

	for i := range nhaccongs {
		if status, found := invitationMap[nhaccongs[i].UserID]; found {
			nhaccongs[i].InvitationStatus = status
		} else {
			nhaccongs[i].InvitationStatus = 4
		}
	}

	return c.JSON(fiber.Map{
		"message":    "success",
		"data":       nhaccongs,
		"user_id":    userLogin.UserID,
		"page":       page,
		"totalItems": totalItems,
	})
}
