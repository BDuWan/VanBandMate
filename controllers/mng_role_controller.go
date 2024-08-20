package controllers

import (
	"fmt"
	"lms/initializers"
	"lms/models"
	"lms/structs"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/gorm"
)

func GetRolePage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetAccRoles")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	return c.Render("pages/management/mng-role/index", fiber.Map{
		"This":           2,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

func APIGetRole(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetRole")
	var roles []models.Role
	DB := initializers.DB

	if err := DB.Where("deleted", false).Find(&roles).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"data": roles,
	})
}

func APIGetRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIGetRoleID")
	var role models.Role
	roleId := c.Params("id")
	var rolePermissions []models.RolePermission
	DB := initializers.DB

	if err := DB.Where("role_id", roleId).Where("deleted", false).First(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("error")
	}
	if err := DB.Where("role_id", roleId).Find(&rolePermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}
	var permissions []int
	for _, rp := range rolePermissions {
		permissions = append(permissions, rp.PermissionID)
	}
	roleData := fiber.Map{
		"name":        role.Name,
		"describe":    role.Describe,
		"permissions": permissions,
	}

	return c.JSON(fiber.Map{
		"data": roleData,
	})
}

func APIPostCreateRole(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostCreateRole")
	var role models.Role
	userLogin := GetSessionUser(c)

	DB := initializers.DB

	roleForm := new(structs.RoleForm)
	if err := c.BodyParser(roleForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON(err.Error())
	}

	if strings.TrimSpace(roleForm.Name) == "" {
		return c.JSON("Tên vai trò không được để trống")
	}

	if err := DB.Where("name", roleForm.Name).First(&models.Role{}).Error; err == nil {
		return c.JSON("Tên vai trò đã được sử dụng")
	}
	if len(roleForm.Permissions) < 1 {
		return c.JSON("Cần chọn ít nhất 1 quyền")
	}

	role.Name = roleForm.Name
	role.Describe = roleForm.Describe
	role.Deleted = false
	role.CreatedBy = userLogin.UserID
	role.DeletedAt = time.Now()
	role.CreatedAt = time.Now()

	if err := DB.Create(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Can not create role")
	}

	for _, permissionId := range roleForm.Permissions {
		var rolePer models.RolePermission
		rolePer.RoleID = role.RoleID
		rolePer.PermissionID = permissionId
		rolePer.CreatedBy = userLogin.UserID
		rolePer.CreatedAt = time.Now()
		if err := DB.Create(&rolePer).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.RedirectBack("")
		}
	}
	return c.JSON("success")
}

func APIPutUpdateRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPutUpdateRoleID")
	var role models.Role
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	roleId := c.Params("id")

	idRole, err := strconv.Atoi(roleId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

	}
	if idRole <= 3 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Role cannot be update")
		return c.JSON("Vai trò này không thể được sửa đổi")
	}
	form := new(structs.RoleForm)
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Dữ liệu không hợp lệ")
	}

	if err := DB.Where("role_id = ?", roleId).Where("deleted", false).First(&role).Error; err != nil {

		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Vai trò đã bị xóa")
	}
	if strings.TrimSpace(form.Name) == "" {
		return c.JSON("Tên vai trò không được để trống")
	}
	if form.Name != role.Name {
		if err := DB.Where("name", form.Name).First(&models.Role{}).Error; err == nil {
			return c.JSON("Tên vai trò đã được sử dụng")
		}
	}

	if len(form.Permissions) < 1 {
		return c.JSON("Cần chọn ít nhất 1 quyền")
	}

	role.Name = form.Name
	role.Describe = form.Describe
	role.UpdatedBy = userLogin.UserID
	role.UpdatedAt = time.Now()

	if err := DB.Model(&role).Updates(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.JSON("Đã xảy ra lỗi khi cập nhật vai trò")
	}

	if err := DB.Where("role_id = ?", roleId).Delete(models.RolePermission{}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("Đã xảy ra lỗi khi cập nhật quyền")
	}

	for _, permissionId := range form.Permissions {
		var rolePer models.RolePermission

		rolePer.RoleID = role.RoleID
		rolePer.PermissionID = permissionId
		rolePer.CreatedBy = userLogin.UserID
		rolePer.CreatedAt = time.Now()

		if err := DB.Create(&rolePer).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
			return c.JSON("Đã xảy ra lỗi khi cập nhật quyền")
		}
	}

	return c.JSON("success")
}

func APIDeleteRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIDeleteRoleID")
	var role models.Role
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	roleId := c.Params("id")

	idRole, err := strconv.Atoi(roleId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

	}
	if idRole <= 3 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Role > 5")

		return c.JSON("Vai trò này không thể xóa")
	}

	if err := DB.Where("role_id", roleId).Where("deleted", false).First(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]:" + err.Error())
		return c.JSON("Không tìm thấy vai trò")
	}
	if role.NumberUser > 0 {
		return c.JSON("Có " + strconv.Itoa(role.NumberUser) + " tài khoản thuộc vai trò này. Vui lòng xóa hoặc chuyển các tài khoản này sang vai trò khác trước khi xóa vai trò")
	}

	role.Deleted = true
	role.DeletedBy = userLogin.UserID
	role.DeletedAt = time.Now()

	if err := DB.Model(&role).Updates(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]:" + err.Error())
		return c.JSON("Can not delete this role.")
	}

	return c.JSON("success")
}

func APIPostCreateRole1(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostCreateRole")
	//var role models.Role
	//DB := initializers.DB
	var roleForm structs.SignUpForm
	if err := c.BodyParser(&roleForm); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + "Format User Fail")
	}
	return c.JSON("Success")
}

func GetAccCreateRole(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetAccCreateRole")
	var permissions []models.Permission

	if err := initializers.DB.Find(&permissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Not found permissions in create role:  " + err.Error())
	}

	return c.Render("pages/accounts/roles/create", fiber.Map{
		"Permissions": permissions,
		"Ctx":         c,
	}, "layouts/main")
}

func APIPostAccRoles(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostAccRoles")
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
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
	}

	return c.JSON(fiber.Map{
		"draw":            req.Draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            roles,
	})
}

func DeleteAccRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: DeleteAccRoleID")
	var role models.Role
	var rolePer models.RolePermission
	//var user models.User1
	userLogin := GetSessionUser(c)
	DB := initializers.DB
	roleId := c.Params("id")

	idRole, err := strconv.Atoi(roleId)
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

	}
	if idRole < 5 {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Role > 5")

		return c.JSON("Can not delete this role.")
	}

	if err := DB.First(&role, roleId).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Not found role ID: " + err.Error())
		return c.JSON("Can not delete this role.")
	}

	role.Deleted = true
	role.DeletedBy = userLogin.UserID
	role.DeletedAt = time.Now()

	if err := DB.Model(&role).Updates(&role).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Can not update role:  " + err.Error())
		return c.JSON("Can not delete this role.")
	}

	if err := DB.Model(&models.RolePermission{}).Where(
		"role_id = ?", roleId).Updates(&rolePer).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Can not delete rol in role permission" + err.Error())
		return c.JSON("Can not delete this role.")
	}

	//user.RoleID = 5
	if err := DB.Model(&models.User1{}).Where("role_id", roleId).Updates(map[string]interface{}{"role_id": gorm.Expr("type_user_id")}).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: Can not update user when update role:  " + err.Error())
		return c.JSON("Can not update user when delete this role.")
	}

	return c.JSON("Success")
}

func GetAccRoleID(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetAccRoleID")
	var role models.Role
	var perMe []models.RolePermission
	var permissions []models.Permission

	DB := initializers.DB
	roleId := c.Params("id")

	err := DB.Where("deleted", false).First(&role, roleId).Error
	if err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Select("permission_id").Where("role_id = ? and deleted = ?", roleId, false).Find(&perMe).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.RedirectBack("")
	}

	if err := DB.Find(&permissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())

		return c.RedirectBack("")
	}

	return c.Render("pages/accounts/roles/edit", fiber.Map{
		"Role":        role,
		"PerMe":       perMe,
		"Permissions": permissions,
		"Ctx":         c,
	}, "layouts/main")
}
