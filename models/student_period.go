package models

import "time"

type StudentPeriod struct {
	StudentPeriodID int       `json:"student_period_id"`
	StudentID       int       `json:"student_id"`
	PeriodID        int       `json:"period_id"`
	Period          Period    `gorm:"foreignKey:PeriodID;references:PeriodID"`
	Deleted         bool      `json:"deleted"`
	DeletedBy       int       `json:"deleted_by"`
	UpdatedBy       int       `json:"updated_by"`
	CreatedBy       int       `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}
