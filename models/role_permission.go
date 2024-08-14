package models

import (
	"encoding/gob"
	"time"
)

func init() {
	// Đăng ký loại dữ liệu User với gob
	gob.Register(RolePermission{})
	gob.Register([]RolePermission{})
}

type RolePermission struct {
	RolePermissionID int        `json:"role_permission_id" gorm:"primaryKey;autoIncrement"`
	RoleID           int        `json:"role_id"`
	Role             Role       `gorm:"foreignKey:RoleID;references:RoleID"`
	PermissionID     int        `json:"permission_id"`
	Permission       Permission `gorm:"foreignKey:PermissionID;references:PermissionID"`
	CreatedBy        int        `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
}
