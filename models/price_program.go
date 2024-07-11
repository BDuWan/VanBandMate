package models

import "time"

type PriceProgram struct {
	PriceProgramID int       `json:"price_program_id" gorm:"primaryKey;autoIncrement"`
	Price          int       `json:"price"`
	Commission     int       `json:"commission"`
	Description    string    `json:"description" gorm:"size:255"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Deleted        bool      `json:"deleted"`
	Default        bool      `json:"default"`
	DeletedBy      int       `json:"deleted_by"`
	UpdatedBy      int       `json:"updated_by"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
