package controllers

import (
	"fmt"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"strconv"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
)

func GetMngSaleBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngSaleBusiness")
	return c.Render("pages/managements/sale_business/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func APIPostMngSaleBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostMngSaleBusiness")
	var users []models.User
	var query *gorm.DB
	DB := initializers.DB
	userLogin := GetSessionUser(c)

	switch userLogin.TypeUserID {
	case 1:
		{

			// query = DB.Model(&models.User{}).Joins("CommissionUser").Where(
			// 	"users.type_user_id = ? or users.type_user_id = ?", 2, 3).Where(
			// 	"users.deleted", false).Where(
			// 	"users.verify", true).Where(
			// 	"users.state", true)
			var commission_users []models.CommissionUser
			query = DB.Model(&models.CommissionUser{}).Joins("User").Where(
				"User.type_user_id = ? or User.type_user_id = ?", 2, 3).Where(
				"commission_users.deleted", false).Where(
				"User.verify", true).Where(
				"User.state", true)
			if err := query.Find(&commission_users).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			}
			return c.JSON(fiber.Map{
				"data": commission_users,
			})
		}
	case 2, 3:
		{
			query = DB.Model(&models.User{}).Joins("TypeUser").Where(
				"users.type_user_id = ?", userLogin.TypeUserID).Where(
				"users.deleted", false).Where(
				"users.verify", true).Where(
				"users.state", true)
		}
	default:
		{
			query = DB.Model(&models.User{}).Joins("TypeUser").Where(
				"users.referral_code = ?", userLogin.CodeUser).Where(
				"users.deleted", false).Where(
				"users.verify", true).Where(
				"users.state", true)
		}
	}
	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

func APIPostStudentMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostStudentMngSaleBusinessID")
	var users []models.User
	var saleBusiness models.User
	var query *gorm.DB
	DB := initializers.DB
	userId := c.Params("id")

	if err := DB.Where("deleted", false).First(&saleBusiness, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	query = DB.Where("referral_code", saleBusiness.CodeUser).Where("deleted", false).Where("type_user_id", 5)

	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

func APIPostPeriodMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPeriodMngSaleBusinessID")
	var periods []models.Period
	//var saleBusiness models.User
	var query *gorm.DB
	DB := initializers.DB
	userId := c.Params("id")

	// if err := DB.Where("deleted", false).First(&saleBusiness, userId).Error; err != nil {
	// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	// }

	query = DB.Where("user_id", userId).Where("deleted", false)

	if err := query.Find(&periods).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": periods,
	})
}

func APIPostInstructorMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostInstructorMngSaleBusinessID")
	var users []models.User
	var saleBusiness models.User
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	userId := c.Params("id")

	if err := DB.Where("deleted", false).First(&saleBusiness, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	query = DB.Where("referral_code", saleBusiness.CodeUser).Where("deleted", false).Where("type_user_id", 4)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&users).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "users.email"
		}
	case 2:
		{
			sortColumn = "users.first_name"
		}
	case 3:
		{
			sortColumn = "users.last_name"
		}
	case 4:
		{
			sortColumn = "users.phone_number"
		}
	case 5:
		{
			sortColumn = "users.address"
		}
	default:
		{
			sortColumn = "users.email"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("users.email LIKE ? "+
			"or users.address LIKE ? "+
			"or users.phone_number LIKE ? "+
			"or users.last_name LIKE ? "+
			"or users.first_name LIKE ? ", search, search, search, search, search)
	}

	query.Find(&users).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            users,
	})
}

func GetMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngSaleBusinessID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.Redirect("/errors/403")
	}
	var user models.User
	var commissionUser models.CommissionUser
	userId := c.Params("id")

	DB := initializers.DB

	if err := DB.Model(&models.CommissionUser{}).Joins("User").Joins("Period").Where(
		"User.deleted", false).Where("User.user_id", userId).First(&commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		//return c.RedirectBack("")
	}

	if err := DB.Where("deleted", false).First(&user, userId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		//return c.RedirectBack("")
	}


	return c.Render("pages/managements/sale_business/detail", fiber.Map{
		"CommissionUser": commissionUser,
		"User":           user,
		"Ctx":            c,
	}, "layouts/main")
}

func APIPostPaymentsMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentsMngSaleBusinessID")
	var commissionUsers []models.CommissionUser
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userId := c.Params("id")

	query = DB.Model(&models.CommissionUser{}).Where("commission_users.user_id", userId).Where("commission_users.deleted", false)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&commissionUsers).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "commission_users.commission"
		}
	case 2:
		{
			sortColumn = "commission_users.commission_total"
		}
	case 3:
		{
			sortColumn = "commission_users.commission_paid"
		}
	case 4:
		{
			sortColumn = "commission_users.commission_debt"
		}
	default:
		{
			sortColumn = "commission_users.commission"
			sortDir = "desc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("commission_users.commission LIKE ? "+
			"or commission_users.commission_total LIKE ? "+
			"or commission_users.commission_debt LIKE ? "+
			"or commission_users.commission_paid LIKE ? ", search, search, search, search)
	}

	query.Find(&commissionUsers).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&commissionUsers).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            commissionUsers,
	})
}

func APIPostHistoryPayMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostHistoryPayMngSaleBusinessID")
	var historyCommissions []models.HistoryCommission
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userId := c.Params("id")

	query = DB.Model(&models.HistoryCommission{}).Where(
		"history_commissions.user_id", userId).Where(
		"history_commissions.type", "pay").Where(
		"history_commissions.deleted", false)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&historyCommissions).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "history_commissions.commission_total"
		}
	case 2:
		{
			sortColumn = "history_commissions.commission_paid"
		}
	case 3:
		{
			sortColumn = "history_commissions.commission_debt"
		}
	case 4:
		{
			sortColumn = "history_commissions.created_at"
		}
	default:
		{
			sortColumn = "history_commissions.created_at"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("history_commissions.commission_total LIKE ? "+
			"or history_commissions.commission_paid LIKE ? "+
			"or history_commissions.commission_debt LIKE ? "+
			"or history_commissions.created_at LIKE ? ", search, search, search, search)
	}

	query.Find(&historyCommissions).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&historyCommissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            historyCommissions,
	})
}

func APIPostHistoryGainMngSaleBusinessID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentsMngSaleBusinessID")
	var historyCommissions []models.HistoryCommission
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB
	userId := c.Params("id")

	query = DB.Model(&models.HistoryCommission{}).Where(
		"history_commissions.user_id", userId).Where(
		"history_commissions.type", "gain").Where(
		"history_commissions.deleted", false)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var totalRecords int64
	var filteredRecords int64

	var sortColumn string
	var sortDir string

	sortDir = req.Order[0].Dir
	query.Find(&historyCommissions).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "history_commissions.commission_total"
		}
	case 2:
		{
			sortColumn = "history_commissions.commission_paid"
		}
	case 3:
		{
			sortColumn = "history_commissions.commission_debt"
		}
	case 4:
		{
			sortColumn = "history_commissions.created_at"
		}
	default:
		{
			sortColumn = "history_commissions.created_at"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("history_commissions.commission_total LIKE ? "+
			"or history_commissions.commission_paid LIKE ? "+
			"or history_commissions.commission_debt LIKE ? "+
			"or history_commissions.created_at LIKE ? ", search, search, search, search)
	}

	query.Find(&historyCommissions).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&historyCommissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            historyCommissions,
	})
}
func GetMngCreatePaymentSaleBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngCreatePaymentSaleBusiness")
	var pricePro models.PriceUser

	if err := initializers.DB.Where("deleted", false).Where("user_id", c.Params("id")).First(&pricePro).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}
	return c.Render("pages/managements/sale_business/payments/create", fiber.Map{
		"PriceProgramID": pricePro.PriceProgramID,
		"UserID":         c.Params("id"),
		"Ctx":            c,
	}, "layouts/main")

}

func PostMngEditPaymentSaleBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostMngCreatePaymentSaleBusiness")
	var commissionUser models.CommissionUser
	var history models.HistoryCommission
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.CommissionEdit)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if form.Commission == 0 {
		return c.JSON("Commission has not been entered")
	}
	if err := DB.Where("user_id", form.UserID).First(&commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found payment ID: " + err.Error())
		return c.JSON("Can not edit commission (not found)")
	}
	commissionUser.CreatedBy = userLogin.UserID

	if err := DB.Model(&commissionUser).Where("user_id", form.UserID).Updates(commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not edit commission" + err.Error())
	}
	history.UserID = form.UserID
	history.CommissionTotal = commissionUser.CommissionTotal
	history.CommissionDebt = commissionUser.CommissionDebt
	history.CommissionPaid = commissionUser.CommissionPaid
	history.Description = "Admin update commission: " + strconv.Itoa(form.Commission)
	history.Deleted = false
	history.CreatedBy = userLogin.UserID
	history.CreatedAt = time.Now()
	if err := DB.Table("history_commissions").Create(&history).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not create history" + err.Error())
	}
	return c.JSON("Success")

}

func PostMngCreatePaymentSaleBusiness(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostMngCreatePaymentSaleBusiness")
	var commissionUser models.CommissionUser
	var history models.HistoryCommission
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.CommissionPay)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Not true format data")
	}

	if form.Price == 0 {
		return c.JSON("Commission has not been entered")
	}
	if err := DB.Where("user_id", form.UserID).First(&commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found payment ID: " + err.Error())
		return c.JSON("Can not pay (not found)")
	}
	if form.Price > commissionUser.CommissionDebt {
		return c.JSON("The payment amount cannot exceed the outstanding balance")
	}
	formatPrice := FormatFloat64(form.Price)
	commissionUser.CommissionPaid = FormatFloat64(commissionUser.CommissionPaid + formatPrice)
	commissionUser.CommissionDebt = FormatFloat64(commissionUser.CommissionDebt - formatPrice)
	commissionUser.CreatedBy = userLogin.UserID
	if err := DB.Model(&commissionUser).Where("user_id", form.UserID).Updates(commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not pay" + err.Error())
	}

	if commissionUser.CommissionDebt == 0 {
		if err := DB.Model(&commissionUser).Where("user_id = ?", form.UserID).Update("commission_debt", 0).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return c.JSON("Can not pay" + err.Error())
		}
	}
	history.UserID = form.UserID
	history.CommissionTotal = commissionUser.CommissionTotal
	history.CommissionDebt = FormatFloat64(commissionUser.CommissionDebt)
	history.CommissionPaid = FormatFloat64(commissionUser.CommissionPaid)
	history.Description = "Admin pay: " + strconv.FormatFloat(form.Price, 'f', 3, 64)
	history.Type = "pay"
	history.Deleted = false
	history.CreatedBy = userLogin.UserID
	history.CreatedAt = time.Now()
	history.DeletedAt = time.Now()
	if err := DB.Table("history_commissions").Create(&history).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON("Can not create history" + err.Error())
	}
	return c.JSON("Success")
}

// func GainCommissionSaleBusiness(referralCode string, timeCreate time.Time) string {
// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GainCommissionSaleBusiness")
// 	var commissionUser models.CommissionUser
// 	var prices []models.PriceProgram
// 	var history models.HistoryCommission
// 	var defaultPrice models.PriceProgram

// 	DB := initializers.DB

// 	if err := DB.Model(&models.CommissionUser{}).Joins("User").Where(
// 		"User.deleted", false).Where(
// 		"User.code_user", referralCode).First(&commissionUser).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		return "Not found sale or business with your referral code"
// 	}

// 	var userId = commissionUser.UserID

// 	if err := DB.Where("deleted", false).Where("default", false).Find(&prices).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		return "Cannot find price"
// 	}

// 	var useDefaultPrice = true
// 	var commission int
// 	for _, price := range prices {
// 		if timeCreate.After(price.StartTime) && timeCreate.Before(price.EndTime) {
// 			commission = price.Commission
// 			useDefaultPrice = false
// 			break
// 		}
// 	}
// 	if useDefaultPrice {
// 		if err := DB.Where("deleted", false).Where("default", true).First(&defaultPrice).Error; err != nil {
// 			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 			return "Cannot find default price"
// 		}
// 		commission = defaultPrice.Commission
// 	}

// 	commissionUser.CommissionTotal += commission
// 	commissionUser.CommissionDebt += commission
// 	//commission_user.CreatedBy = user_id
// 	if err := DB.Model(&commissionUser).Where("user_id", userId).Updates(commissionUser).Error; err != nil {
// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
// 		return "Cannot gain commission for sale or business"
// 	}

//		history.UserID = commissionUser.User.UserID
//		history.CommissionTotal = commissionUser.CommissionTotal
//		history.CommissionDebt = commissionUser.CommissionDebt
//		history.CommissionPaid = commissionUser.CommissionPaid
//		history.Description = "Student sign up, gain: " + strconv.Itoa(commission)
//		history.Type = "gain"
//		history.Deleted = false
//		//history.CreatedBy = user_id
//		history.CreatedAt = time.Now()
//		if err := DB.Table("history_commissions").Create(&history).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//			return "Can not create history"
//		}
//		return "ok"
//	}
func GainNumberStudentSaleBusiness(user_id int, referralCode string, timeCreate time.Time) string {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GainNumberStudentSaleBusiness")
	var commissionUser models.CommissionUser

	DB := initializers.DB

	if err := DB.Model(&models.CommissionUser{}).Joins("User").Joins("Period").Where(
		"User.deleted", false).Where(
		"User.code_user", referralCode).First(&commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Not found sale or business with this referral code"
	}

	var period_id = commissionUser.PeriodID

	var commissionUser_id = commissionUser.CommissionUserID
	if commissionUser.PeriodID == 0 || commissionUser.Period.PeriodEnd.Before(timeCreate) {
		var config_payment models.ConfigPayment
		if err := DB.First(&config_payment).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return "Not found config payment"
		}
		var period models.Period
		period.NumberStudent = 1
		period.StudentPaid = 0
		period.UserID = commissionUser.UserID
		period.PeriodStart = time.Now()
		period.PeriodEnd = AddWorkingDays(time.Now(), config_payment.NumberDayPeriod)
		period.CommissionDefault = FormatFloat64(config_payment.CommissionDefault)
		period.CommissionBonus = FormatFloat64(config_payment.CommissionBonus)
		period.NumberToGetBonus = config_payment.NumberToGetBonus
		period.Deleted = false
		period.CreatedAt = time.Now()
		period.CreatedAt = time.Now()
		period.DeletedAt = time.Now()
		period.UpdatedAt = time.Now()
		if err := DB.Table("periods").Create(&period).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
			return "Can not create new period"
		}
		period_id = int(period.PeriodID)
		commissionUser.PeriodID = int(period.PeriodID)
		commissionUser.UpdatedAt = time.Now()
		if err := DB.Model(&commissionUser).Where("commission_user_id = ?", commissionUser_id).Updates(commissionUser).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update commissionUser:  " + err.Error())
			return "Error updating commission user: " + err.Error()
		}

	} else {
		commissionUser.Period.NumberStudent++
		if err := DB.Model(&commissionUser.Period).Updates(commissionUser.Period).Error; err != nil {
			return "Error updating period" + err.Error()
		}
	}
	var student_period models.StudentPeriod
	student_period.StudentID = user_id
	student_period.PeriodID = period_id
	student_period.Deleted = false
	student_period.CreatedAt = time.Now()
	student_period.DeletedAt = time.Now()
	student_period.UpdatedAt = time.Now()
	if err := DB.Table("student_periods").Create(&student_period).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Can not create new student period"
	}
	return "ok"
}

func GainCommissionSaleBusiness(referralCode string, periodID int) string {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GainCommissionSaleBusiness")
	var commissionUser models.CommissionUser
	var period models.Period
	//var prices []models.PriceProgram
	//var history models.HistoryCommission
	//var defaultPrice models.PriceProgram

	DB := initializers.DB

	if err := DB.Model(&models.CommissionUser{}).Joins("User").Where(
		"User.deleted", false).Where(
		"User.code_user", referralCode).First(&commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Not found sale or business with this referral code"
	}
	if err := DB.Where("period_id", periodID).First(&period).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return "Not found any period"
	}
	var commission_user_id = commissionUser.CommissionUserID

	period.StudentPaid += 1
	if period.StudentPaid > period.NumberToGetBonus {
		period.Commission += FormatFloat64(period.CommissionDefault + period.CommissionBonus)
		commissionUser.CommissionTotal += FormatFloat64(period.CommissionDefault + period.CommissionBonus)
		commissionUser.CommissionDebt += FormatFloat64(period.CommissionDefault + period.CommissionBonus)
	} else {
		period.Commission += FormatFloat64(period.CommissionDefault)
		commissionUser.CommissionTotal += FormatFloat64(period.CommissionDefault)
		commissionUser.CommissionDebt += FormatFloat64(period.CommissionDefault)
	}

	if err := DB.Model(&period).Updates(period).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update period:  " + err.Error())
		return "Can not update this period."
	}
	if err := DB.Model(&commissionUser).Where("commission_user_id", commission_user_id).Updates(commissionUser).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update commissionUser:  " + err.Error())
		return "Can not update this commissionUser."
	}

	//var userId = commissionUser.UserID

	// if err := DB.Where("deleted", false).Where("default", false).Find(&prices).Error; err != nil {
	// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	// 	return "Cannot find price"
	// }

	// var useDefaultPrice = true
	// var commission int
	// for _, price := range prices {
	// 	if timeCreate.After(price.StartTime) && timeCreate.Before(price.EndTime) {
	// 		commission = price.Commission
	// 		useDefaultPrice = false
	// 		break
	// 	}
	// }
	// if useDefaultPrice {
	// 	if err := DB.Where("deleted", false).Where("default", true).First(&defaultPrice).Error; err != nil {
	// 		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	// 		return "Cannot find default price"
	// 	}
	// 	commission = defaultPrice.Commission
	// }

	// commissionUser.CommissionTotal += commission
	// commissionUser.CommissionDebt += commission
	// //commission_user.CreatedBy = user_id
	// if err := DB.Model(&commissionUser).Where("user_id", userId).Updates(commissionUser).Error; err != nil {
	// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	// 	return "Cannot gain commission for sale or business"
	// }

	// history.UserID = commissionUser.User.UserID
	// history.CommissionTotal = commissionUser.CommissionTotal
	// history.CommissionDebt = commissionUser.CommissionDebt
	// history.CommissionPaid = commissionUser.CommissionPaid
	// history.Description = "Student sign up, gain: " + strconv.Itoa(commission)
	// history.Type = "gain"
	// history.Deleted = false
	// //history.CreatedBy = user_id
	// history.CreatedAt = time.Now()
	// if err := DB.Table("history_commissions").Create(&history).Error; err != nil {
	// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	// 	return "Can not create history"
	// }
	return "ok"
}

func AddWorkingDays(startDate time.Time, workingDays int) time.Time {
	daysAdded := 0
	currentDate := startDate

	for daysAdded < workingDays {
		currentDate = currentDate.AddDate(0, 0, 1)
		if currentDate.Weekday() != time.Sunday {
			daysAdded++
		}
	}
	return currentDate
}

func DeleteMngSaleBusinessPaymentID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteMngSaleBusinessPaymentID")
	if GetSessionUser(c).TypeUserID > 1 {
		return c.JSON("Permission Denied")
	}
	var payment models.Payment

	userLogin := GetSessionUser(c)
	DB := initializers.DB
	paymentId := c.Params("id")

	if err := DB.First(&payment, paymentId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found payment ID: " + err.Error())
		return c.JSON("Can not delete this payment.")
	}

	payment.Deleted = true
	payment.DeletedBy = userLogin.UserID
	payment.DeletedAt = time.Now()

	if err := DB.Model(&payment).Updates(payment).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update payment:  " + err.Error())
		return c.JSON("Can not delete this payment.")
	}

	return c.JSON("Success")
}

func DailyUpdatePeriod() {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Daily Update Period Active")
	DB := initializers.DB
	var commissionUsers []models.CommissionUser
	if err := DB.Joins("Period").Find(&commissionUsers).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Daily Update Period: Cannot find commission user " + err.Error())
		return
	}
	for _, commissionUser := range commissionUsers {
		commissionUser_Id := commissionUser.CommissionUserID
		if commissionUser.PeriodID == 0 || commissionUser.Period.PeriodEnd.Before(time.Now()) {
			var period models.Period
			period.NumberStudent = 0
			period.StudentPaid = 0
			period.UserID = commissionUser.UserID
			period.PeriodStart = time.Now()
			period.PeriodEnd = AddWorkingDays(time.Now(), 50)
			period.Deleted = false
			period.CreatedAt = time.Now()
			period.CreatedAt = time.Now()
			period.DeletedAt = time.Now()
			period.UpdatedAt = time.Now()
			if err := DB.Table("periods").Create(&period).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Daily Update Period: Cannot creating new period " + err.Error())
				return
			}

			commissionUser.PeriodID = int(period.PeriodID)
			commissionUser.UpdatedAt = time.Now()
			if err := DB.Model(&commissionUser).Where("commission_user_id = ?", commissionUser_Id).Updates(commissionUser).Error; err != nil {
				outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Daily Update Period: Can not update commissionUser:  " + err.Error())
				return
			}

		}
	}
}
