package models

import (
	"encoding/gob"
	"fmt"
	"time"
)

func init() {
	// Đăng ký loại dữ liệu User với gob
	gob.Register(User{})
}

type User struct {
	UserID                 int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	TypeUserID             int       `json:"type_user_id"`
	TypeUser               TypeUser  `gorm:"foreignKey:TypeUserID;references:TypeUserID"`
	CodeUser               string    `json:"code_user" gorm:"size:20"`
	FirstName              string    `json:"first_name" gorm:"size:50"`
	LastName               string    `json:"last_name" gorm:"size:50"`
	NameBusiness           string    `json:"name_business" gorm:"size:255"`
	FullNameRepresentative string    `json:"full_name_representative" gorm:"size:255"`
	RoleID                 int       `json:"role_id"`
	Role                   Role      `gorm:"foreignKey:RoleID;references:RoleID"`
	Email                  string    `json:"email" validate:"required,min=11,max=35" gorm:"size:50;uniqueIndex"`
	PhoneNumber            string    `json:"phone_number" gorm:"size:15"`
	Address                string    `json:"address" gorm:"size:255"`
	Username               string    `json:"username" validate:"required,min=8,max=20" gorm:"size:20;uniqueIndex"`
	Password               string    `json:"password" validate:"required,min=8,max=20" gorm:"size:255"`
	Token                  string    `json:"token" gorm:"size:255"`
	ReferralCode           string    `json:"referral_code" gorm:"size:20"`
	Session                string    `json:"session" gorm:"size:100"`
	State                  bool      `json:"state"`
	Paid                   bool      `json:"paid"`
	Verify                 bool      `json:"verify"`
	Deleted                bool      `json:"deleted"`
	DeletedBy              int       `json:"deleted_by"`
	UpdatedBy              int       `json:"updated_by"`
	CreatedBy              int       `json:"created_by"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at"`
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
