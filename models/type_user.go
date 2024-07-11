package models

type TypeUser struct {
	TypeUserID int    `json:"type_user_id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"size:100"`
}
