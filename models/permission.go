package models

type Permission struct {
	PermissionID int    `json:"permission_id" gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name" gorm:"size:50"`
	Permission   string `json:"permission" gorm:"size:50"`
	Href         string `json:"href" gorm:"size:50"`
	Icon         string `json:"icon" gorm:"size:50"`
	Level        int    `json:"level" gorm:"size:50"`
}
