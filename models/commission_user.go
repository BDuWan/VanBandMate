package models

import "time"

type CommissionUser struct {
	CommissionUserID int       `json:"commission_user_id"`
	UserID           int       `json:"user_id"`
	User             User      `gorm:"foreignKey:UserID;references:UserID"`
	CommissionTotal  float64   `json:"commission_total"`
	CommissionPaid   float64   `json:"commission_paid"`
	CommissionDebt   float64   `json:"commission_debt"`
	PeriodID         int       `json:"period_id"`
	Period           Period    `gorm:"foreignKey:PeriodID;references:PeriodID"`
	Deleted          bool      `json:"deleted"`
	DeletedBy        int       `json:"deleted_by"`
	UpdatedBy        int       `json:"updated_by"`
	CreatedBy        int       `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}
