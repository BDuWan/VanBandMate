package controllers

//
//import (
//	"fmt"
//	"lms/initializers"
//	"lms/models"
//	"lms/structs"
//	"os"
//	"time"
//
//	"github.com/go-playground/validator/v10"
//	"github.com/gofiber/fiber/v2"
//	"github.com/zetamatta/go-outputdebug"
//	"gorm.io/gorm"
//)
//
//func GetMngPayments(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngPayments")
//	var user models.User
//	var commissionUser models.CommissionUser
//	userId := GetSessionUser(c).UserID
//
//	DB := initializers.DB
//
//	if err := DB.Model(&models.CommissionUser{}).Joins("User").Joins("Period").Where(
//		"User.deleted", false).Where("User.user_id", userId).First(&commissionUser).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		//return c.RedirectBack("")
//	}
//
//	if err := DB.Where("deleted", false).First(&user, userId).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//	//return c.Re
//	return c.Render("pages/managements/payments/index", fiber.Map{
//		"TypeUserID":     GetSessionUser(c).TypeUserID,
//		"CommissionUser": commissionUser,
//		"User":           user,
//		"Ctx":            c,
//	}, "layouts/main")
//}
//
//func APIPostPaymentPrices(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentPrices")
//	var prices []models.PriceProgram
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Where("deleted", false)
//
//	if err := c.BodyParser(&req); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Cannot parse JSON",
//		})
//	}
//
//	var totalRecords int64
//	var filteredRecords int64
//
//	var sortColumn string
//	var sortDir string
//
//	sortDir = req.Order[0].Dir
//	query.Find(&prices).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "price"
//		}
//	case 2:
//		{
//			sortColumn = "commission"
//		}
//	case 3:
//		{
//			sortColumn = "description"
//		}
//	default:
//		{
//			sortColumn = "start_time"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("price LIKE ? "+
//			"or commission LIKE ? "+
//			"or description LIKE ? ", search, search, search)
//	}
//
//	query.Find(&prices).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&prices).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            prices,
//	})
//}
//
//func DeletePaymentPriceID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeletePaymentPriceID")
//	var price models.PriceProgram
//	DB := initializers.DB
//	priceId := c.Params("id")
//
//	if err := DB.First(&price, priceId).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found price ID: " + err.Error())
//		return c.JSON("Can not delete this price.")
//	}
//
//	price.Deleted = true
//	price.DeletedBy = GetSessionUser(c).UserID
//	price.DeletedAt = time.Now()
//
//	if err := DB.Model(&price).Updates(&price).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update price when delete:  " + err.Error())
//		return c.JSON("Can not delete this price.")
//	}
//
//	return c.JSON("Success")
//}
//
//func GetPaymentCreatePrice(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetPaymentCreatePrice")
//	return c.Render("pages/managements/payments/prices/create", fiber.Map{
//		"Ctx": c,
//	}, "layouts/main")
//}
//
//func PostPaymentCreatePrice(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostPaymentCreatePrice")
//	var price models.PriceProgram
//	userLogin := GetSessionUser(c)
//
//	DB := initializers.DB
//
//	form := new(structs.PriceProgram)
//	if err := c.BodyParser(form); err != nil {
//
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON(err.Error())
//	}
//
//	if err := validator.New().Struct(form); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Validate create price: " + err.Error())
//		return c.JSON("Invalid input: " + err.Error())
//	}
//
//	if form.Price == 0 {
//		return c.JSON("Price can not be blank")
//	}
//	if form.Commission == 0 {
//		return c.JSON("Commission can not be blank")
//	}
//	if form.Commission >= form.Price {
//		return c.JSON("The commission must be less than the program price")
//	}
//	if form.StartTime == "" {
//		return c.JSON("Haven't chosen a start time")
//	}
//	if form.EndTime == "" {
//		return c.JSON("Haven't chosen a end time")
//	}
//	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
//	if errStr != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
//	}
//	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
//	if errEnd != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
//	}
//	checkTime := CheckValidTime(startTime, endTime)
//	if checkTime != "ok" {
//		return c.JSON(checkTime)
//	}
//	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
//	if errOff != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
//	}
//
//	price.Price = form.Price
//	price.Commission = form.Commission
//	price.Description = form.Description
//	price.StartTime = startTime.Add(offset)
//	price.EndTime = endTime.Add(offset)
//	price.Deleted = false
//	price.Default = false
//	price.CreatedBy = userLogin.UserID
//	price.DeletedAt = time.Now()
//	price.CreatedAt = time.Now()
//
//	if err := DB.Create(&price).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not create price")
//	}
//
//	return c.JSON("Success")
//}
//
//func APIPostPaymentProSales(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentProSales")
//	var prices []models.CommissionUser
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	// query = DB.Select("a.price_program_id", "b.price_user_id", "a.commission", "a.price", "c.user_id", "c.first_name", "c.last_name", "c.name_business", "c.full_name_representative").Table(
//	// 	"price_programs a").Joins(
//	// 	"LEFT JOIN price_users b on a.price_program_id = b.price_program_id").Joins(
//	// 	"RIGHT JOIN users c on b.user_id = c.user_id").Where(
//	// 	"c.deleted", false).Where(
//	// 	"c.type_user_id < 4").Where(
//	// 	"c.type_user_id > 1").Where(
//	// 	"b.deleted", false)
//
//	query = DB.Model(&models.CommissionUser{}).Joins("User").Where(
//		"commission_users.deleted", false).Where(
//		"User.deleted", false)
//	if err := c.BodyParser(&req); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Cannot parse JSON",
//		})
//	}
//
//	var totalRecords int64
//	var filteredRecords int64
//
//	var sortColumn string
//	var sortDir string
//
//	sortDir = req.Order[0].Dir
//	query.Find(&prices).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "User.first_name"
//		}
//	case 2:
//		{
//			sortColumn = "User.last_name"
//		}
//	case 3:
//		{
//			sortColumn = "User.name_business"
//		}
//	case 4:
//		{
//			sortColumn = "User.full_name_representative"
//		}
//	case 5:
//		{
//			sortColumn = "commission_total"
//		}
//	case 6:
//		{
//			sortColumn = "commission_paid"
//		}
//	case 7:
//		{
//			sortColumn = "commission_debt"
//		}
//	default:
//		{
//			sortColumn = "User.last_name"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("User.first_name LIKE ? "+
//			"or User.last_name LIKE ? "+
//			"or User.name_business LIKE ? "+
//			"or User.full_name_representative LIKE ? "+
//			"or commission_total LIKE ? "+
//			"or commission_paid LIKE ? "+
//			"or commission_debt LIKE ? ", search, search, search, search, search, search, search)
//	}
//
//	query.Find(&prices).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&prices).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            prices,
//	})
//}
//
//func GetPaymentCreatePriceSaleID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetPaymentCreatePriceSaleID")
//	var prices []models.PriceProgram
//	var priceUser models.PriceUser
//	DB := initializers.DB
//
//	userId := c.Params("uid")
//	if err := DB.Where("deleted", false).Find(&prices).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	if err := DB.Where("deleted", false).First(&priceUser, c.Params("id")).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.Render("pages/managements/payments/sales/create", fiber.Map{
//		"PriceUser": priceUser,
//		"Prices":    prices,
//		"UserID":    userId,
//		"Ctx":       c,
//	}, "layouts/main")
//}
//
//func PostPaymentCreatePriceSaleID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostPaymentCreatePriceSaleID")
//	var priceUser models.PriceUser
//	userLogin := GetSessionUser(c)
//
//	DB := initializers.DB
//
//	form := new(structs.PriceProgramUser)
//	if err := c.BodyParser(form); err != nil {
//
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON(err.Error())
//	}
//
//	if err := DB.Where("deleted", false).First(&priceUser, form.PriceUserID).Error; err != nil {
//
//		var priceU models.PriceUser
//		priceU.PriceProgramID = form.PriceProgramID
//		priceU.UserID = form.UserID
//		priceU.Deleted = false
//		priceU.DeletedAt = time.Now()
//		priceU.UpdatedAt = time.Now()
//		priceU.CreatedAt = time.Now()
//		priceU.CreatedBy = userLogin.UserID
//		if err := DB.Create(&priceU).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//			return c.JSON("Can not create price")
//		}
//		return c.JSON("Success")
//	}
//	priceUser.Deleted = true
//	priceUser.DeletedBy = userLogin.UserID
//	priceUser.DeletedAt = time.Now()
//
//	if err := DB.Model(&priceUser).Updates(priceUser).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON("Can not create price")
//	}
//
//	var priceU models.PriceUser
//	priceU.PriceProgramID = form.PriceProgramID
//	priceU.UserID = form.UserID
//	priceU.Deleted = false
//	priceU.DeletedAt = time.Now()
//	priceU.UpdatedAt = time.Now()
//	priceU.CreatedAt = time.Now()
//	priceU.CreatedBy = userLogin.UserID
//	if err := DB.Create(&priceU).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON("Can not create price")
//	}
//	return c.JSON("Success")
//}
//
//func GetPaymentConfig(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetPaymentConfig")
//	var config_payment models.ConfigPayment
//	DB := initializers.DB
//
//	if err := DB.First(&config_payment).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON("Not found config payment")
//	}
//
//	return c.Render("pages/managements/payments/edit", fiber.Map{
//		"ConfigPayment": config_payment,
//		"Ctx":           c,
//	}, "layouts/main")
//}
//
//func UpdateMngPaymentsConfig(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngPaymentsConfig")
//	var config models.ConfigPayment
//	DB := initializers.DB
//	form := new(structs.ConfigPayment)
//	if err := c.BodyParser(form); err != nil {
//
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON(err.Error())
//	}
//
//	if err := DB.First(&config).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON("Not found config payment")
//	}
//
//	if(form.CommissionBonus == 0 || form.CommissionDefault == 0 || form.NumberDay == 0 || form.NumberStudent == 0){
//		return c.JSON("Invalid input")
//	}
//
//	config.CommissionDefault = FormatFloat64(form.CommissionDefault)
//	config.CommissionBonus = FormatFloat64(form.CommissionBonus)
//	config.NumberToGetBonus = form.NumberStudent
//	config.NumberDayPeriod = form.NumberDay
//	config.UpdatedAt = time.Now()
//
//	if err := DB.Model(&config).Updates(&config).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not update config payment")
//	}
//
//	return c.JSON("Success")
//}
//func GetMngPaymentsPriceID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetMngPaymentsPriceID")
//	var price models.PriceProgram
//
//	if err := initializers.DB.Where("deleted", false).First(&price, c.Params("id")).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//	return c.Render("pages/managements/payments/prices/edit", fiber.Map{
//		"Price": price,
//		"Ctx":   c,
//	}, "layouts/main")
//}
//
//func UpdateMngPaymentsPriceID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateMngPaymentsPriceID")
//	var price models.PriceProgram
//	DB := initializers.DB
//	priceId := c.Params("id")
//	form := new(structs.PriceProgram)
//	if err := c.BodyParser(form); err != nil {
//
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON(err.Error())
//	}
//
//	if err := DB.Where("deleted", false).First(&price, priceId).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return c.JSON("Not found")
//	}
//
//	if form.Price == 0 {
//		return c.JSON("Price can not be blank")
//	}
//	if form.Commission == 0 {
//		return c.JSON("Commission can not be blank")
//	}
//	if form.Commission >= form.Price {
//		return c.JSON("The commission must be less than the program price")
//	}
//	if form.StartTime == "" {
//		return c.JSON("Haven't chosen a start time")
//	}
//	if form.EndTime == "" {
//		return c.JSON("Haven't chosen a end time")
//	}
//	startTime, errStr := time.Parse("2006-01-02T15:04", form.StartTime)
//	if errStr != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errStr.Error())
//	}
//	endTime, errEnd := time.Parse("2006-01-02T15:04", form.EndTime)
//	if errEnd != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errEnd.Error())
//	}
//
//	checkTime := CheckValidTimeUpdate(startTime, endTime, priceId)
//	if checkTime != "ok" {
//		return c.JSON(checkTime)
//	}
//
//	offset, errOff := time.ParseDuration(os.Getenv("OFFSET_TIME"))
//	if errOff != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + errOff.Error())
//	}
//
//	price.Price = form.Price
//	price.Commission = form.Commission
//	price.Description = form.Description
//	price.StartTime = startTime.Add(offset)
//	price.EndTime = endTime.Add(offset)
//	price.UpdatedBy = GetSessionUser(c).UserID
//	price.UpdatedAt = time.Now()
//
//	if err := DB.Model(&price).Updates(&price).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not update price")
//	}
//
//	return c.JSON("Success")
//}
//
//func APIPostPaymentSales(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentSales")
//	var commissionUsers []models.CommissionUser
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//	userLogin := GetSessionUser(c)
//
//	query = DB.Model(&models.CommissionUser{}).Where("commission_users.user_id", userLogin.UserID).Where("commission_users.deleted", false)
//
//	if err := c.BodyParser(&req); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Cannot parse JSON",
//		})
//	}
//
//	var totalRecords int64
//	var filteredRecords int64
//
//	var sortColumn string
//	var sortDir string
//
//	sortDir = req.Order[0].Dir
//	query.Find(&commissionUsers).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "commission_users.commission_total"
//		}
//	case 2:
//		{
//			sortColumn = "commission_users.commission_paid"
//		}
//	case 3:
//		{
//			sortColumn = "commission_users.commission_debt"
//		}
//	default:
//		{
//			sortColumn = "commission_users.commission"
//			sortDir = "desc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("commission_users.commission_total LIKE ? "+
//			"or commission_users.commission_debt LIKE ? "+
//			"or commission_users.commission_paid LIKE ? ", search, search, search)
//	}
//
//	query.Find(&commissionUsers).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&commissionUsers).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            commissionUsers,
//	})
//}
//
//func APIPostPaymentStudent(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentStudent")
//	var users []models.User
//	var saleBusiness models.User
//	var query *gorm.DB
//	DB := initializers.DB
//	userId := GetSessionUser(c).UserID
//
//	if err := DB.Where("deleted", false).First(&saleBusiness, userId).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	query = DB.Where("referral_code", saleBusiness.CodeUser).Where("deleted", false).Where("type_user_id", 5)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"data": users,
//	})
//}
//
//func APIPostPaymentPeriod(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentPeriod")
//	var periods []models.Period
//	//var saleBusiness models.User
//	var query *gorm.DB
//	DB := initializers.DB
//	userId := GetSessionUser(c).UserID
//
//	// if err := DB.Where("deleted", false).First(&saleBusiness, userId).Error; err != nil {
//	// 	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	// }
//
//	query = DB.Where("user_id", userId).Where("deleted", false)
//
//	if err := query.Find(&periods).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"data": periods,
//	})
//}
//
//func APIPostPaymentHistoryPay(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentHistoryPay")
//	var historyCommissions []models.HistoryCommission
//	var query *gorm.DB
//	DB := initializers.DB
//	userId := GetSessionUser(c).UserID
//
//	query = DB.Model(&models.HistoryCommission{}).Where(
//		"history_commissions.user_id", userId).Where(
//		"history_commissions.type", "pay").Where(
//		"history_commissions.deleted", false)
//
//	if err := query.Find(&historyCommissions).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"data": historyCommissions,
//	})
//}
//
//func APIPostPaymentHistory(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostPaymentsMngSaleBusinessID")
//	var historyCommissions []models.HistoryCommission
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//	userLogin := GetSessionUser(c)
//
//	query = DB.Model(&models.HistoryCommission{}).Where("history_commissions.user_id", userLogin.UserID).Where("history_commissions.deleted", false)
//
//	if err := c.BodyParser(&req); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Cannot parse JSON",
//		})
//	}
//
//	var totalRecords int64
//	var filteredRecords int64
//
//	var sortColumn string
//	var sortDir string
//
//	sortDir = req.Order[0].Dir
//	query.Find(&historyCommissions).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "history_commissions.commission_total"
//		}
//	case 2:
//		{
//			sortColumn = "history_commissions.commission_paid"
//		}
//	case 3:
//		{
//			sortColumn = "history_commissions.commission_debt"
//		}
//	default:
//		{
//			sortColumn = "history_commissions.created_at"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("payments.total LIKE ? "+
//			"or payments.created_at LIKE ? "+
//			"or PriceProgram.price LIKE ? "+
//			"or PriceProgram.commission LIKE ? ", search, search, search)
//	}
//
//	query.Find(&historyCommissions).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&historyCommissions).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            historyCommissions,
//	})
//}
//
//func CheckValidTime(startTime time.Time, endTime time.Time) string {
//	if startTime.After(endTime) {
//		return "The end time must be after the start time"
//	}
//	DB := initializers.DB
//	var prices []models.PriceProgram
//	if err := DB.Where("deleted", false).Where("default", false).Find(&prices).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return "Cannot find price program"
//	}
//	for _, price := range prices {
//		if startTime.After(price.StartTime) && startTime.Before(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//		if endTime.After(price.StartTime) && endTime.Before(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//		if startTime.Before(price.StartTime) && endTime.After(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//	}
//	return "ok"
//}
//
//func CheckValidTimeUpdate(startTime time.Time, endTime time.Time, priceId string) string {
//	if startTime.After(endTime) {
//		return "The end time must be after the start time"
//	}
//	DB := initializers.DB
//	var prices []models.PriceProgram
//	if err := DB.Where("deleted", false).Where("default", false).Where("price_program_id <> ?", priceId).Find(&prices).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//		return "Cannot find price program"
//	}
//	for _, price := range prices {
//		if startTime.After(price.StartTime) && startTime.Before(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//		if endTime.After(price.StartTime) && endTime.Before(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//		if startTime.Before(price.StartTime) && endTime.After(price.EndTime) {
//			return "The time coincides with the price program: " + price.Description
//		}
//	}
//	return "ok"
//}
