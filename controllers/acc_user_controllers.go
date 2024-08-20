package controllers

//
//import (
//	"fmt"
//	"lms/initializers"
//	"lms/models"
//	"lms/structs"
//	"lms/utils"
//	"regexp"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/gofiber/fiber/v2"
//	"github.com/zetamatta/go-outputdebug"
//	"gorm.io/gorm"
//)
//
//func GetAccUsers(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccUsers")
//	return c.Render("pages/accounts/users/index", fiber.Map{
//		"Ctx": c,
//	}, "layouts/main")
//}
//
//func APIPostAccUserAdmin(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccUserAdmin")
//	var users []models.User1
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Model(&models.User1{}).Joins("Role").Where("users.deleted", false).Where("users.type_user_id", 1)
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
//	query.Find(&users).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "users.email"
//		}
//	case 2:
//		{
//			sortColumn = "users.first_name"
//		}
//	case 3:
//		{
//			sortColumn = "users.last_name"
//		}
//	case 4:
//		{
//			sortColumn = "Role.name"
//		}
//	default:
//		{
//			sortColumn = "users.email"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("users.email LIKE ? "+
//			"or Role.name LIKE ? "+
//			"or users.last_name LIKE ? "+
//			"or users.first_name LIKE ? ", search, search, search, search)
//	}
//
//	query.Find(&users).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            users,
//	})
//}
//
//func APIPostAccUserSales(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccUserSales")
//	var users []models.User1
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Model(&models.User1{}).Joins("Role").Where("users.deleted", false).Where("users.type_user_id", 2)
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
//	query.Find(&users).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "users.email"
//		}
//	case 2:
//		{
//			sortColumn = "users.first_name"
//		}
//	case 3:
//		{
//			sortColumn = "users.last_name"
//		}
//	case 4:
//		{
//			sortColumn = "Role.name"
//		}
//	default:
//		{
//			sortColumn = "users.email"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("users.email LIKE ? "+
//			"or Role.name LIKE ? "+
//			"or users.last_name LIKE ? "+
//			"or users.first_name LIKE ? ", search, search, search, search)
//	}
//
//	query.Find(&users).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            users,
//	})
//}
//
//func APIPostAccUserBusiness(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccUserBusiness")
//	var users []models.User1
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Model(&models.User1{}).Joins("Role").Where("users.deleted", false).Where("users.type_user_id", 3)
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
//	query.Find(&users).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "users.email"
//		}
//	case 2:
//		{
//			sortColumn = "users.name_business"
//		}
//	case 3:
//		{
//			sortColumn = "users.full_name_representative"
//		}
//	case 4:
//		{
//			sortColumn = "Role.name"
//		}
//	default:
//		{
//			sortColumn = "users.email"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("users.email LIKE ? "+
//			"or Role.name LIKE ? "+
//			"or users.full_name_representative LIKE ? "+
//			"or users.name_business LIKE ? ", search, search, search, search)
//	}
//
//	query.Find(&users).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            users,
//	})
//}
//
//func APIPostAccUserInstructors(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccUserInstructors")
//	var users []models.User1
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Model(&models.User1{}).Joins("Role").Where("users.deleted", false).Where("users.type_user_id", 4)
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
//	query.Find(&users).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "users.email"
//		}
//	case 2:
//		{
//			sortColumn = "users.first_name"
//		}
//	case 3:
//		{
//			sortColumn = "users.last_name"
//		}
//	case 4:
//		{
//			sortColumn = "Role.name"
//		}
//	default:
//		{
//			sortColumn = "users.email"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("users.email LIKE ? "+
//			"or Role.name LIKE ? "+
//			"or users.last_name LIKE ? "+
//			"or users.first_name LIKE ? ", search, search, search, search)
//	}
//
//	query.Find(&users).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            users,
//	})
//}
//
//func APIPostAccUserStudents(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccUserStudents")
//	var users []models.User1
//	var query *gorm.DB
//	var req structs.ReqBody
//	DB := initializers.DB
//
//	query = DB.Model(&models.User1{}).Joins("Role").Where("users.deleted", false).Where("users.type_user_id", 5)
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
//	query.Find(&users).Count(&totalRecords)
//
//	switch req.Order[0].Column {
//	case 1:
//		{
//			sortColumn = "users.email"
//		}
//	case 2:
//		{
//			sortColumn = "users.first_name"
//		}
//	case 3:
//		{
//			sortColumn = "users.last_name"
//		}
//	case 4:
//		{
//			sortColumn = "Role.name"
//		}
//	default:
//		{
//			sortColumn = "users.email"
//			sortDir = "asc"
//		}
//	}
//	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
//	query = query.Order(orderBy)
//
//	if req.Search.Value != "" {
//		search := "%" + req.Search.Value + "%"
//		query = query.Where("users.email LIKE ? "+
//			"or Role.name LIKE ? "+
//			"or users.last_name LIKE ? "+
//			"or users.first_name LIKE ? ", search, search, search, search)
//	}
//
//	query.Find(&users).Count(&filteredRecords)
//
//	query = query.Offset(req.Start).Limit(req.Length)
//
//	if err := query.Find(&users).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"draw":            req.Draw,
//		"recordsTotal":    totalRecords,
//		"recordsFiltered": filteredRecords,
//		"data":            users,
//	})
//}
//
//func GetAccCreateUser(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccCreateUser")
//	var roles []models.Role
//	var typeUsers []models.TypeUser
//
//	if err := initializers.DB.Find(&roles).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found roles in create user:  " + err.Error())
//	}
//
//	if err := initializers.DB.Find(&typeUsers).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found type users in create user:  " + err.Error())
//	}
//
//	return c.Render("pages/accounts/users/create", fiber.Map{
//		"TypeUsers": typeUsers,
//		"Roles":     roles,
//		"Ctx":       c,
//	}, "layouts/main")
//}
//
//func PostAccCreateUser(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostAccCreateUser")
//	var user structs.AccUser
//	var account models.User1
//	var codeUser string
//	var validator = "ok"
//
//	DB := initializers.DB
//	if err := c.BodyParser(&user); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User1 Fail")
//	}
//
//	userLogin := GetSessionUser(c)
//
//	if user.TypeUserID != 3 {
//		validator = ValidatorCreateAccInput(user)
//	}
//
//	if validator != "ok" {
//		return c.JSON(validator)
//	}
//
//	if err := DB.Where("email", user.Email).First(&models.User1{}).Error; err != nil {
//		account.TypeUserID = user.TypeUserID
//		account.FirstName = user.FirstName
//		account.LastName = user.LastName
//		account.RoleID = user.RoleID
//		account.Email = user.Email
//		account.PhoneNumber = user.PhoneNumber
//		account.Address = user.Address
//		account.Username = user.Email
//		account.Password = utils.HashingPassword(user.Password)
//		account.ReferralCode = user.ReferralCode
//		account.Session = ""
//		account.NameBusiness = user.NameBusiness
//		account.FullNameRepresentative = user.FullNameRepresentative
//		account.State = true
//		account.Verify = false
//		account.Deleted = false
//		account.CreatedAt = time.Now()
//		account.DeletedAt = time.Now()
//		account.UpdatedAt = time.Now()
//		account.CreatedBy = userLogin.UserID
//
//		var checkReferralCode models.User1
//		if user.ReferralCode != "" {
//			if err := DB.Where("type_user_id = ? or type_user_id = ?", 2, 3).Where(
//				"code_user", user.ReferralCode).First(&checkReferralCode).Error; err != nil {
//				return c.JSON("Please enter the correct referral code or leave it blank")
//			}
//		}
//
//		if err := DB.Create(&account).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
//			return c.JSON("Can not create account")
//		}
//
//		switch user.TypeUserID {
//		case 1:
//			{
//				codeUser = "ADM"
//			}
//		case 2:
//			{
//				codeUser = "SAL"
//			}
//		case 3:
//			{
//				codeUser = "BUSIN"
//			}
//		case 4:
//			{
//				codeUser = "INSTR"
//			}
//		case 5:
//			{
//				codeUser = "STUD"
//			}
//
//		}
//
//		account.CodeUser = codeUser + fmt.Sprintf("%04d", account.UserID)
//
//		if err := DB.Updates(&account).Error; err != nil {
//			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create account")
//			return c.JSON("Can not create account")
//		}
//
//		//if account.TypeUserID == 2 || account.TypeUserID == 3 {
//		//	var commissionUser models.CommissionUser
//		//	commissionUser.UserID = account.UserID
//		//	commissionUser.CommissionTotal = 0
//		//	commissionUser.CommissionPaid = 0
//		//	commissionUser.CommissionDebt = 0
//		//	commissionUser.PeriodID = 0
//		//	commissionUser.Deleted = false
//		//	commissionUser.CreatedAt = time.Now()
//		//	commissionUser.CreatedBy = account.UserID
//		//	if err := DB.Create(&commissionUser).Error; err != nil {
//		//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not create commission user")
//		//	}
//		//}
//		return c.JSON("Success")
//	}
//
//	return c.JSON("Email already exists")
//}
//
//func DeleteAccUserID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteAccUserID")
//	//var user models.User1
//	DB := initializers.DB
//	userId := c.Params("id")
//
//	idUser, err := strconv.Atoi(userId)
//	if err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//	}
//	if idUser < 2 {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: User1 > 1")
//
//		return c.JSON("Can not delete this user.")
//	}
//
//	//if err := DB.First(&user, userId).Error; err != nil {
//	//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found user ID: " + err.Error())
//	//	return c.JSON("Can not delete this user.")
//	//}
//	//
//	//user.State = false
//	//user.Verify = false
//	//user.Deleted = true
//	//user.DeletedBy = GetSessionUser(c).UserID
//	//user.DeletedAt = time.Now()
//
//	if err := DB.Where("user_id", userId).Delete(&models.User1{}).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update user when delete:  " + err.Error())
//		return c.JSON("Can not delete this user.")
//	}
//
//	return c.JSON("Success")
//}
//
//func GetAccUserID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccUserID")
//	var user models.User1
//	var roles []models.Role
//	var typeUsers []models.TypeUser
//
//	DB := initializers.DB
//	userId := c.Params("id")
//	err := DB.Where("deleted", false).First(&user, userId).Error
//	if err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Where("deleted", false).Find(&roles).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.RedirectBack("")
//	}
//
//	if err := DB.Select("type_user_id", "name").Find(&typeUsers).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
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
//	return c.Render("pages/accounts/users/edit", fiber.Map{
//		"User1": user,
//		//"data":       data,
//		"Roles": roles,
//		//"RoleID":     user.RoleID,
//		"TypeUsers": typeUsers,
//		//"TypeUserID": user.TypeUserID,
//		"Ctx": c,
//	}, "layouts/main")
//}
//
//func UpdateAccUserID(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateAccUserID")
//	var user structs.AccUser
//	var account models.User1
//	DB := initializers.DB
//
//	userId := c.Params("id")
//
//	if err := c.BodyParser(&user); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User1 Fail")
//	}
//
//	if err := DB.Where("user_id", userId).First(&account).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Format User1 Fail")
//		return c.JSON("Can not found account")
//	}
//
//	account.TypeUserID = user.TypeUserID
//	account.FirstName = user.FirstName
//	account.LastName = user.LastName
//	account.RoleID = user.RoleID
//	account.Email = user.Email
//	account.PhoneNumber = user.PhoneNumber
//	account.Address = user.Address
//	account.ReferralCode = user.ReferralCode
//	account.NameBusiness = user.NameBusiness
//	account.FullNameRepresentative = user.FullNameRepresentative
//	account.UpdatedAt = time.Now()
//	account.UpdatedBy = GetSessionUser(c).UserID
//
//	if err := DB.Model(&account).Updates(&account).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + "Can not update account")
//		return c.JSON("Can not update account")
//	}
//
//	return c.JSON("Success")
//
//}
//
//func UpdateAccStateUser(c *fiber.Ctx) error {
//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateAccStateUser")
//	var user models.User1
//	DB := initializers.DB
//	form := new(structs.FormStateUser)
//	if err := c.BodyParser(form); err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not update state")
//	}
//	if form.ID < 2 {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: user id < 2")
//
//		return c.JSON("Can not update state")
//	}
//	if err := DB.Where("user_id", form.ID).First(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not update state")
//	}
//	user.State = form.State
//	user.Session = ""
//	if err := DB.Save(&user).Error; err != nil {
//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
//
//		return c.JSON("Can not update state")
//	}
//	return c.JSON("Success")
//}
//
//func ValidatorCreateAccInput(user structs.AccUser) string {
//	if strings.TrimSpace(user.FirstName) == "" {
//		return "Firstname can not be blank"
//	}
//	if strings.TrimSpace(user.LastName) == "" {
//		return "Lastname can not be blank"
//	}
//	if strings.TrimSpace(user.Email) == "" {
//		return "Email can not be blank"
//	}
//	if strings.TrimSpace(user.PhoneNumber) == "" {
//		return "PhoneNumber can not be blank"
//	}
//	if strings.TrimSpace(user.Address) == "" {
//		return "Address can not be blank"
//	}
//	//if user.Username == "" {
//	//	return "Username can not be blank"
//	//}
//	//if user.Password == "" {
//	//	return "Password can not be blank"
//	//}
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
//	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
//	regexE := regexp.MustCompile(regexEmail)
//	if !regexE.MatchString(user.Email) {
//		return "Invalid Email"
//	}
//
//	regexPhone := "^[0-9]{7,}$"
//	regexP := regexp.MustCompile(regexPhone)
//	if !regexP.MatchString(user.PhoneNumber) {
//		return "Invalid Phone number"
//	}
//
//	//regexUsername := "^[a-zA-Z0-9]{8,}$"
//	//regexU := regexp.MustCompile(regexUsername)
//	//if !regexU.MatchString(user.Username) {
//	//	return "Username must have at least 8 characters, including only letters and numbers"
//	//}
//
//	testPassword := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[!@#$%^&*()?]"}
//	for _, test := range testPassword {
//		t, _ := regexp.MatchString(test, user.Password)
//		if !t {
//			return "Password too weak"
//		}
//	}
//	return "ok"
//}
