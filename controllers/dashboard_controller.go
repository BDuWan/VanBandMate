package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"time"
	"vanbandmate/initializers"
	"vanbandmate/models"
)

func GetDashBoardPage(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: GetMngUserPage")
	userLogin := GetSessionUser(c)
	sess, _ := SessAuth.Get(c)
	permissions := sess.Get("rolePermission")
	var allPermissions []models.Permission
	var provinces []models.Province
	DB := initializers.DB

	if err := DB.Find(&allPermissions).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}

	if err := DB.Find(&provinces).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.Redirect("/errors/404")
	}

	var years []int
	if err := DB.Model(&models.Contract{}).
		Where("date IS NOT NULL").
		Where("deleted", false).
		Select("DISTINCT EXTRACT(YEAR FROM date) AS year").
		Order("year ASC").
		Pluck("year", &years).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: " + err.Error())
		return c.JSON("error")
	}

	return c.Render("pages/dashboard/index", fiber.Map{
		"This":           5,
		"Permissions":    permissions,
		"AllPermissions": allPermissions,
		"Provinces":      provinces,
		"Years":          years,
		"User":           userLogin,
		"Ctx":            c,
	}, "layouts/main")
}

func APIPostDashboardFilter(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM]: APIPostUserFilter")
	var users []models.User
	DB := initializers.DB

	form := new(struct {
		Year  int `json:"year"`
		Month int `json:"month"`
	})
	if err := c.BodyParser(form); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM1]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	var chuloadaiStats []struct {
		UserID        uint
		ContractCount int
		TotalPrice    int64
	}

	rootQuery := DB.Table("contracts")
	if form.Year != 0 {
		rootQuery = rootQuery.Where("YEAR(contracts.date) = ?", form.Year)
	}

	if form.Month != 0 {
		rootQuery = rootQuery.Where("MONTH(contracts.date) = ?", form.Month)
	}

	rootQuery.
		Select("chuloadai_id as user_id, COUNT(*) as contract_count, SUM(price) as total_price").
		Group("chuloadai_id").
		Scan(&chuloadaiStats)

	var nhaccongStats []struct {
		UserID        uint
		ContractCount int
		TotalPrice    int64
	}
	rootQuery.
		Select("nhaccong_id as user_id, COUNT(*) as contract_count, SUM(price) as total_price").
		Group("nhaccong_id").
		Scan(&nhaccongStats)

	stats := map[uint]struct {
		ContractCount int
		TotalPrice    int64
	}{}
	for _, stat := range chuloadaiStats {
		if existing, found := stats[stat.UserID]; found {
			stats[stat.UserID] = struct {
				ContractCount int
				TotalPrice    int64
			}{
				ContractCount: existing.ContractCount + stat.ContractCount,
				TotalPrice:    existing.TotalPrice + stat.TotalPrice,
			}
		} else {
			stats[stat.UserID] = struct {
				ContractCount int
				TotalPrice    int64
			}{
				ContractCount: stat.ContractCount,
				TotalPrice:    stat.TotalPrice,
			}
		}
	}
	for _, stat := range nhaccongStats {
		if existing, found := stats[stat.UserID]; found {
			stats[stat.UserID] = struct {
				ContractCount int
				TotalPrice    int64
			}{
				ContractCount: existing.ContractCount + stat.ContractCount,
				TotalPrice:    existing.TotalPrice + stat.TotalPrice,
			}
		} else {
			stats[stat.UserID] = struct {
				ContractCount int
				TotalPrice    int64
			}{
				ContractCount: stat.ContractCount,
				TotalPrice:    stat.TotalPrice,
			}
		}
	}

	var userIDs []uint
	for userID := range stats {
		userIDs = append(userIDs, userID)
	}

	query := DB.Model(&models.User{}).
		Joins("Role").Joins("Province").Joins("District").Joins("Ward").
		Where("users.deleted = ?", false).
		Where("users.user_id IN ?", userIDs)

	if err := query.Find(&users).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [VBM2]: " + err.Error())
		return c.JSON(fiber.Map{
			"message": "Đã xảy ra lỗi khi lấy dữ liệu",
		})
	}

	for i, user := range users {
		if stat, exists := stats[uint(user.UserID)]; exists {
			users[i].CountContract = stat.ContractCount // Gán số lượng hợp đồng
			users[i].SumPrice = int(stat.TotalPrice)    // Gán tổng giá trị hợp đồng
		}
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}
