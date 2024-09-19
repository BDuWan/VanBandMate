package initializers

import (
	"github.com/zetamatta/go-outputdebug"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
	"vanbandmate/models"
	"vanbandmate/utils"
)

var DB *gorm.DB

//func ConnectToDatabase() {
//	var err error
//	dsn := os.Getenv("USER_DB") + ":" + os.Getenv("PASSWORD") + `@tcp(127.0.0.1:3306)/` + os.Getenv("DATABASE") + `?charset=utf8mb4&parseTime=True&loc=Local`
//
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		//Logger: logger.Default.LogMode(logger.Info),
//		Logger: logger.Default.LogMode(logger.Silent),
//	})
//
//	if err != nil {
//		outputdebug.String("[LMS]: " + err.Error())
//
//	}
//}

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("USER_DB") + ":" + os.Getenv("PASSWORD") + `@tcp(us-cluster-east-01.k8s.cleardb.net:3306)/` + os.Getenv("DATABASE") + `?charset=utf8mb4&parseTime=True&loc=Local`

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		outputdebug.String("[LMS]: " + err.Error())

	}
}

func MigrateDB() {
	if err := DB.AutoMigrate(&models.Role{}); err != nil {
		outputdebug.String("[LMS]: " + err.Error())
	}
	if err := DB.AutoMigrate(&models.Permission{}); err != nil {
		outputdebug.String("[LMS]: " + err.Error())
	}
	if err := DB.AutoMigrate(&models.RolePermission{}); err != nil {
		outputdebug.String("[LMS]: " + err.Error())
	}
	//if err := DB.AutoMigrate(&models.TypeUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		outputdebug.String("[LMS]: " + err.Error())
	}
	//if err := DB.AutoMigrate(&models.StudyProgram{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.Course{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.StudyProgramUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.CourseUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.Document{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.Assignment{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.AssignmentUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.Lesson{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.FileUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.PriceProgram{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.Payment{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
	//if err := DB.AutoMigrate(&models.PriceUser{}); err != nil {
	//	outputdebug.String("[LMS]: " + err.Error())
	//}
}
func GenData() {
	roles := []models.Role{
		{Name: "Administrator"},
		{Name: "Sale Agent"},
		{Name: "Business"},
		{Name: "Instructor"},
		{Name: "Student"},
	}

	//typeUsers := []models.TypeUser{
	//	{Name: "Administrator"},
	//	{Name: "Sale Agent"},
	//	{Name: "Business"},
	//	{Name: "Instructor"},
	//	{Name: "Student"},
	//}

	permissions := []models.Permission{
		{Name: "Home", Permission: "home"},
		{Name: "Dashboard", Permission: "dashboard"},
		{Name: "StudyPrograms", Permission: "study_programs"},
		{Name: "Add StudyProgram", Permission: "add_study_program"},
		{Name: "Courses", Permission: "courses"},
		{Name: "Add Course", Permission: "add_course"},
		{Name: "Management Instructors", Permission: "mng_instructors"},
		{Name: "Management Students", Permission: "mng_students"},
		{Name: "Management Sales Agents and Business", Permission: "mng_sale_business"},
		{Name: "Management Courses", Permission: "mng_courses"},
		{Name: "Management StudyPrograms", Permission: "mng_study_programs"},
		{Name: "Management Users", Permission: "account_users"},
		{Name: "Management Roles", Permission: "account_roles"},
		{Name: "Management Information", Permission: "management_info"},
		{Name: "Management Account", Permission: "management_account"},
		{Name: "Management Payments", Permission: "mng_payments"},
	}

	rolePermissions := []models.RolePermission{
		{RoleID: 1, PermissionID: 1},
		{RoleID: 1, PermissionID: 2},
		{RoleID: 1, PermissionID: 7},
		{RoleID: 1, PermissionID: 8},
		{RoleID: 1, PermissionID: 9},
		{RoleID: 1, PermissionID: 10},
		{RoleID: 1, PermissionID: 11},
		{RoleID: 1, PermissionID: 12},
		{RoleID: 1, PermissionID: 13},
		{RoleID: 1, PermissionID: 14},
		{RoleID: 1, PermissionID: 15},
		{RoleID: 1, PermissionID: 16},

		{RoleID: 2, PermissionID: 1},
		{RoleID: 2, PermissionID: 2},
		{RoleID: 2, PermissionID: 7},
		{RoleID: 2, PermissionID: 8},
		{RoleID: 2, PermissionID: 14},
		{RoleID: 2, PermissionID: 16},

		{RoleID: 3, PermissionID: 1},
		{RoleID: 3, PermissionID: 2},
		{RoleID: 3, PermissionID: 7},
		{RoleID: 3, PermissionID: 8},
		{RoleID: 3, PermissionID: 14},
		{RoleID: 3, PermissionID: 16},

		{RoleID: 4, PermissionID: 1},
		{RoleID: 4, PermissionID: 3},
		{RoleID: 4, PermissionID: 5},

		{RoleID: 5, PermissionID: 1},
		{RoleID: 5, PermissionID: 3},
		{RoleID: 5, PermissionID: 4},
		{RoleID: 5, PermissionID: 5},
		{RoleID: 5, PermissionID: 6},
	}

	user := models.User{
		FirstName:   "Administrators",
		RoleID:      1,
		Email:       "Admin@gmail.com",
		PhoneNumber: "0123456789",
		Password:    utils.HashingPassword("Admin@#$2024"),
		State:       true,
		Verify:      true,
		Deleted:     false,
		UpdatedBy:   1,
		CreatedBy:   1,
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
		DeletedAt:   time.Now(),
	}

	//priceProgram := models.PriceProgram{
	//	Price:       10000,
	//	Commission:  2000,
	//	Description: "Default",
	//	Deleted:     false,
	//	CreatedBy:   1,
	//	CreatedAt:   time.Now(),
	//	UpdatedAt:   time.Now(),
	//}
	//
	//for _, data := range typeUsers {
	//	if err := DB.Create(&data).Error; err != nil {
	//		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	//	}
	//}

	for _, data := range roles {
		if err := DB.Create(&data).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	for _, data := range permissions {
		if err := DB.Create(&data).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	for _, data := range rolePermissions {
		if err := DB.Create(&data).Error; err != nil {
			outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
		}
	}

	if err := DB.Create(&user).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	//if err := DB.Create(&priceProgram).Error; err != nil {
	//	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	//}

}
