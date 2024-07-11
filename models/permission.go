package models

type Permission struct {
	PermissionID int    `json:"permission_id" gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name" gorm:"size:50"`
	Permission   string `json:"permission" gorm:"size:50"`
}
