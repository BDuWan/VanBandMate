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
		return c.Redirect("/errors/404")
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
		Joins("JOIN provinces ON users.province_code = provinces.code").
		Joins("JOIN districts ON users.district_code = districts.code").
		Joins("JOIN wards ON users.ward_code = wards.code").
		Joins("JOIN roles ON users.role_id = roles.role_id").
		Joins("JOIN role_permissions ON roles.role_id = role_permissions.role_id").
		Where("users.deleted = ?", false).
		Where("role_permissions.permission_id = ?", 10)

	//query := DB.Model(&models.User{}).
	//	Joins("Province").Joins("District").Joins("Ward").
	//	Joins("Role").Joins("Role__RolePermission").
	//	Where("users.deleted", false).
	//	Where("RolePermission.permission_id", 10)

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
		"message":      "success",
		"data":         nhaccongs,
		"user_id":      userLogin.UserID,
		"page":         page,
		"hiringNewsId": hiringNews.HiringNewsID,
		"totalItems":   totalItems,
	})
}

func PostHiringInvite(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostHiringInvite")
	type Request struct {
		HiringNewsID int `json:"hiringNewsID"`
		NhaccongID   int `json:"nhaccongID"`
	}

	var req Request

	// Parse dữ liệu JSON từ request body
	if err := c.BodyParser(&req); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM1]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	hiringNewsID := req.HiringNewsID
	nhaccongID := req.NhaccongID

	DB := initializers.DB
	var hiringNews models.HiringNews
	var invitations []models.Invitation

	if err := DB.Where("hiring_news_id", hiringNewsID).
		Where("deleted", false).
		First(&hiringNews).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM2]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Tin tuyển dụng đã bị xóa",
		})
	}

	if err := DB.Where("hiring_news_id", hiringNewsID).
		Where("nhaccong_id", nhaccongID).
		Find(&invitations).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM3]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	if len(invitations) > 0 {
		invitations[0].Status = 0

		if err := DB.Save(&invitations[0]).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM4]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi khi cập nhật",
			})
		}
	} else {
		newInvitation := models.Invitation{
			HiringNewsID: hiringNewsID,
			NhaccongID:   nhaccongID,
			Status:       0,
			InviteAt:     time.Now(),
		}

		if err := DB.Create(&newInvitation).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM5]: " + err.Error())
			return c.JSON(fiber.Map{
				"icon":    "error",
				"message": "Đã xảy ra lỗi khi gửi lời mời",
			})
		}
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Gửi lời mời thành công",
	})
}

func PostHiringCancelInvite(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostHiringCancelInvite")
	type Request struct {
		HiringNewsID int `json:"hiringNewsID"`
		NhaccongID   int `json:"nhaccongID"`
	}

	var req Request

	// Parse dữ liệu JSON từ request body
	if err := c.BodyParser(&req); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	hiringNewsID := req.HiringNewsID
	nhaccongID := req.NhaccongID
	DB := initializers.DB
	//var userHiringNewsList []models.UserHiringNews

	if err := DB.Model(&models.Invitation{}).
		Where("hiring_news_id", hiringNewsID).
		Where("nhaccong_id", nhaccongID).
		Update("status", 2).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Thu hồi lời mời thành công",
	})
}

func GetReceiveInvitationPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetAccRoles")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	return c.Render("pages/invitations/receive-inv", fiber.Map{
		"This":           10,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

func APIPostReceivedInvFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostReceivedInvFilter")
	DB := initializers.DB

	form := new(structs.FormReceivedInvFind)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	page := form.Page
	limit := form.ItemsPerPage

	offset := (page - 1) * limit

	userLogin := GetSessionUser(c)
	query := DB.Model(&models.Invitation{}).
		Joins("HiringNews").Joins("HiringNews.User").
		Joins("HiringNews.Province").Joins("HiringNews.District").Joins("HiringNews.Ward").
		Where("nhaccong_id", userLogin.UserID).
		Where("HiringNews.deleted", false).
		Where("HiringNews__User.deleted", false).
		Where("status IN ?", []int{0, 1, 3})

	var totalItems int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&totalItems).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var invitations []models.Invitation
	if err := query.Offset(offset).Limit(limit).Find(&invitations).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	//Kiểm tra lời mời trùng thời gian
	var contracts []models.Contract
	if err := DB.Where("nhaccong_id", userLogin.UserID).
		Where("deleted", false).Find(&contracts).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi")
	}

	for i := range invitations {
		for _, contract := range contracts {
			if invitations[i].HiringNews.Date.Equal(contract.Date) {
				if invitations[i].Status != 1 {
					invitations[i].Status = 5
					break
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":    "success",
		"data":       invitations,
		"page":       page,
		"totalItems": totalItems,
	})
}

func PostReceivedInvAccept(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: PostNewsApply")
	type Request struct {
		ID uint `json:"id"`
	}
	var req Request

	// Parse dữ liệu JSON từ request body
	if err := c.BodyParser(&req); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Đã xảy ra lỗi",
		})
	}

	userLogin := GetSessionUser(c)

	invitationID := req.ID
	DB := initializers.DB
	var invitation models.Invitation

	if err := DB.Where("invitation_id", invitationID).
		Joins("HiringNews").Joins("HiringNews.User").
		Joins("HiringNews.Province").Joins("HiringNews.District").Joins("HiringNews.Ward").
		Where("HiringNews.deleted", false).
		Where("HiringNews__User.deleted", false).
		First(&invitation).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Lời mời đã bị xóa",
		})
	}

	if err := DB.Model(&models.Invitation{}).
		Where("hiring_news_id", invitation.HiringNewsID).
		Update("status", 2).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không thể cập nhật trạng thái lời mời",
		})
	}

	invitation.Status = 1
	if err := DB.Save(&invitation).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không thể cập nhật trạng thái lời mời",
		})
	}

	contract := models.Contract{
		ChuloadaiID:   invitation.HiringNews.ChuloadaiID,
		NhaccongID:    invitation.NhaccongID,
		ProvinceCode:  invitation.HiringNews.ProvinceCode,
		DistrictCode:  invitation.HiringNews.DistrictCode,
		WardCode:      invitation.HiringNews.WardCode,
		AddressDetail: invitation.HiringNews.AddressDetail,
		Status:        1,
		Price:         invitation.HiringNews.Price,
		Date:          invitation.HiringNews.Date,
		Deleted:       false,
		CreatedBy:     userLogin.UserID,
		CreatedAt:     time.Now(),
	}

	if err := DB.Create(&contract).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(fiber.Map{
			"icon":    "error",
			"message": "Không thể tạo hợp đồng",
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(fiber.Map{
		"icon":    "success",
		"message": "Bạn đã chấp nhận lời mời, hợp đồng đã được tạo",
	})
}
