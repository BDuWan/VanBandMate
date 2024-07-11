package models

import "time"

type PriceUser struct {
	PriceUserID    int          `json:"price_user_id" gorm:"primaryKey;autoIncrement"`
	PriceProgramID int          `json:"price_program_id"`
	PriceProgram   PriceProgram `gorm:"foreignKey:PriceProgramID;references:PriceProgramID"`
	UserID         int          `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;references:UserID"`
	Deleted        bool         `json:"deleted"`
	DeletedBy      int          `json:"deleted_by"`
	UpdatedBy      int          `json:"updated_by"`
	CreatedBy      int          `json:"created_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      time.Time    `json:"deleted_at"`
}
