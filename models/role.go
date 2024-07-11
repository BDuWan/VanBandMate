package models

import (
	"time"
)

type Role struct {
	RoleID    int       `json:"role_id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100"`
	Deleted   bool      `json:"deleted"`
	DeletedBy int       `json:"deleted_by"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
