package controllers

import (
	"fmt"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
)

func GetAccRoles(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccRoles")
	return c.Render("pages/accounts/roles/index", fiber.Map{
		"Ctx": c,
	}, "layouts/main")
}

func GetAccCreateRole(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccCreateRole")
	var permissions []models.Permission

	if err := initializers.DB.Find(&permissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found permissions in create role:  " + err.Error())
	}

	return c.Render("pages/accounts/roles/create", fiber.Map{
		"Permissions": permissions,
		"Ctx":         c,
	}, "layouts/main")
}

func APIPostAccRoles(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIPostAccRoles")
	var roles []models.Role
	var query *gorm.DB
	var req structs.ReqBody
	DB := initializers.DB

	query = DB.Where("deleted", false)

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
	query.Find(&roles).Count(&totalRecords)

	switch req.Order[0].Column {
	case 1:
		{
			sortColumn = "roles.name"
		}
	default:
		{
			sortColumn = "roles.name"
			sortDir = "asc"
		}
	}
	orderBy := fmt.Sprintf("%s %s", sortColumn, sortDir)
	query = query.Order(orderBy)

	if req.Search.Value != "" {
		search := "%" + req.Search.Value + "%"
		query = query.Where("roles.name LIKE ? ", search)
	}

	query.Find(&roles).Count(&filteredRecords)

	query = query.Offset(req.Start).Limit(req.Length)

	if err := query.Find(&roles).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            roles,
	})
}

func DeleteAccRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: DeleteAccRoleID")
	var role models.Role
	var rolePer models.RolePermission
	//var user models.User1
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	roleId := c.Params("id")

	idRole, err := strconv.Atoi(roleId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

	}
	if idRole < 5 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Role > 5")

		return c.JSON("Can not delete this role.")
	}

	if err := DB.First(&role, roleId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Not found role ID: " + err.Error())
		return c.JSON("Can not delete this role.")
	}

	role.Deleted = true
	role.DeletedBy = userLogin.UserID
	role.DeletedAt = time.Now()

	if err := DB.Model(&role).Updates(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update role:  " + err.Error())
		return c.JSON("Can not delete this role.")
	}

	rolePer.Deleted = true
	rolePer.DeletedBy = userLogin.UserID
	rolePer.DeletedAt = time.Now()
	if err := DB.Model(&models.RolePermission{}).Where(
		"role_id = ?", roleId).Updates(&rolePer).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not delete rol in role permission" + err.Error())
		return c.JSON("Can not delete this role.")
	}

	//user.RoleID = 5
	if err := DB.Model(&models.User1{}).Where("role_id", roleId).Updates(map[string]interface{}{"role_id": gorm.Expr("type_user_id")}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Can not update user when update role:  " + err.Error())
		return c.JSON("Can not update user when delete this role.")
	}

	return c.JSON("Success")
}

func PostAccCreateRole(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: PostAccCreateRole")
	var role models.Role
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	form := new(structs.RoleForm)
	if err := c.BodyParser(form); err != nil {

		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		return c.JSON(err.Error())
	}

	if err := validator.New().Struct(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Validate create role: " + err.Error())
		return c.JSON("Invalid input: " + err.Error())
	}

	role.Name = form.Name
	role.Deleted = false
	role.CreatedBy = userLogin.UserID
	role.DeletedAt = time.Now()
	role.CreatedAt = time.Now()

	if err := DB.Create(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.JSON("Can not create role")
	}

	if len(form.PermissionID) < 1 {
		return c.JSON("Role must have at least 1 permission")
	}

	for _, permissionId := range form.PermissionID {
		var rolePer models.RolePermission
		rolePer.RoleID = role.RoleID
		rolePer.PermissionID = permissionId
		rolePer.Deleted = false
		rolePer.CreatedBy = userLogin.UserID
		rolePer.CreatedAt = time.Now()
		rolePer.DeletedAt = time.Now()
		if err := DB.Create(&rolePer).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

			return c.RedirectBack("")
		}
	}

	//return c.Redirect("/accounts/roles")
	return c.JSON("Create new role success")
}

func GetAccRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: GetAccRoleID")
	var role models.Role
	var perMe []models.RolePermission
	var permissions []models.Permission

	DB := initializers.DB
	roleId := c.Params("id")

	err := DB.Where("deleted", false).First(&role, roleId).Error
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Select("permission_id").Where("role_id = ? and deleted = ?", roleId, false).Find(&perMe).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Find(&permissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	return c.Render("pages/accounts/roles/edit", fiber.Map{
		"Role":        role,
		"PerMe":       perMe,
		"Permissions": permissions,
		"Ctx":         c,
	}, "layouts/main")
}

func UpdateAccRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: UpdateAccRoleID")
	var role models.Role
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	roleId := c.Params("id")

	idRole, err := strconv.Atoi(roleId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

	}
	if idRole < 5 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Role > 5")

		return c.RedirectBack("")
	}
	form := new(structs.RoleForm)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Where("role_id = ?", roleId).First(&role).Error; err != nil {

		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	role.Name = form.Name
	role.UpdatedBy = userLogin.UserID
	role.UpdatedAt = time.Now()

	if err := DB.Model(&role).Updates(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	if len(form.PermissionID) < 1 {
		return c.RedirectBack("")
	}

	if err := DB.Where("role_id = ?", roleId).Delete(models.RolePermission{}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

		return c.RedirectBack("")
	}

	for _, permissionId := range form.PermissionID {
		var rolePer models.RolePermission

		rolePer.RoleID = role.RoleID
		rolePer.PermissionID = permissionId
		rolePer.Deleted = false
		rolePer.CreatedAt = time.Now()
		rolePer.UpdatedAt = time.Now()
		rolePer.DeletedAt = time.Now()

		if err := DB.Create(&rolePer).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())

			return c.RedirectBack("")
		}
	}

	return c.Redirect("/accounts/roles")
}
