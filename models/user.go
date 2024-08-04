package models

import (
	"fmt"
	"time"
)

//func init() {
//	gob.Register(User{})
//}

type User1 struct {
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

type User struct {
	UserID        int       `gorm:"column:user_id;primaryKey;autoIncrement"`
	FirstName     string    `gorm:"column:first_name;size:50"`
	LastName      string    `gorm:"column:last_name;size:50"`
	RoleID        int       `gorm:"column:role_id"`
	Role          Role      `gorm:"foreignKey:RoleID;references:RoleID"`
	Email         string    `gorm:"column:email;size:50;unique"`
	Gender        int       `gorm:"column:gender"`
	LinkFacebook  string    `json:"link_facebook" gorm:"size:500"`
	ProvinceCode  string    `gorm:"column:province_code;size:50"`
	Province      Province  `gorm:"foreignKey:ProvinceCode;references:Code"`
	DistrictCode  string    `gorm:"column:district_code;size:50"`
	District      District  `gorm:"foreignKey:DistrictCode;references:Code"`
	WardCode      string    `gorm:"column:ward_code;size:50"`
	Ward          Ward      `gorm:"foreignKey:WardCode;references:Code"`
	AddressDetail string    `gorm:"column:address_detail;size:255"`
	PhoneNumber   string    `gorm:"column:phone_number;size:15"`
	Image         string    `gorm:"column:image;size:50"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Password      string    `gorm:"column:password;size:255"`
	Token         string    `gorm:"column:token;size:255"`
	Session       string    `gorm:"column:session;size:100"`
	State         bool      `gorm:"column:state"`
	Verify        bool      `gorm:"column:verify;default:false"`
	Deleted       bool      `gorm:"column:deleted"`
	DeletedBy     int       `gorm:"column:deleted_by"`
	UpdatedBy     int       `gorm:"column:updated_by"`
	CreatedBy     int       `gorm:"column:created_by"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at"`
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
