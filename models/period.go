package models

import (
	"time"
)

type Period struct {
	PeriodID          int64     `json:"period_id" gorm:"primaryKey;autoIncrement"`
	UserID            int       `json:"user_id"`
	User              User      `gorm:"foreignKey:UserID;references:UserID"`
	Commission        float64   `json:"commission"`
	NumberStudent     int       `json:"number_student"`
	StudentPaid       int       `json:"student_paid"`
	PeriodStart       time.Time `json:"period_start"`
	PeriodEnd         time.Time `json:"period_end"`
	CommissionDefault float64   `json:"commission_default"`
	CommissionBonus   float64   `json:"commission_bonus"`
	NumberToGetBonus  int       `json:"number_to_get_bonus"`
	Deleted           bool      `json:"deleted" gorm:"default:false"`
	DeletedBy         int       `json:"deleted_by"`
	UpdatedBy         int       `json:"updated_by"`
	CreatedBy         uint64    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
}
