package models

import "time"

type ClassUser struct {
	ClassUserID int       `json:"class_user_id" gorm:"primaryKey;autoIncrement"`
	ClassID     int       `json:"class_id"`
	Class       Class     `gorm:"foreignKey:ClassID;references:ClassID"`
	UserID      int       `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;references:UserID"`
	Deleted     bool      `json:"deleted"`
	DeletedBy   int       `json:"deleted_by"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
