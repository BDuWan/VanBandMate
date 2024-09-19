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
	UserID           int       `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	FirstName        string    `json:"first_name" gorm:"column:first_name;size:50"`
	LastName         string    `json:"last_name" gorm:"column:last_name;size:50"`
	RoleID           int       `json:"role_id" gorm:"column:role_id"`
	Role             Role      `json:"role" gorm:"foreignKey:RoleID;references:RoleID"`
	Email            string    `json:"email" gorm:"column:email;size:50;unique"`
	Gender           int       `json:"gender" gorm:"column:gender"`
	LinkFacebook     string    `json:"link_facebook" gorm:"column:link_facebook;size:500"`
	ProvinceCode     string    `json:"province_code" gorm:"column:province_code;size:50"`
	Province         Province  `json:"province" gorm:"foreignKey:ProvinceCode;references:Code"`
	DistrictCode     string    `json:"district_code" gorm:"column:district_code;size:50"`
	District         District  `json:"district" gorm:"foreignKey:DistrictCode;references:Code"`
	WardCode         string    `json:"ward_code" gorm:"column:ward_code;size:50"`
	Ward             Ward      `json:"ward" gorm:"foreignKey:WardCode;references:Code"`
	AddressDetail    string    `json:"address_detail" gorm:"column:address_detail;size:255"`
	PhoneNumber      string    `json:"phone_number" gorm:"column:phone_number;size:15"`
	Image            string    `json:"image" gorm:"column:image;size:50"`
	DateOfBirth      time.Time `json:"date_of_birth" gorm:"column:date_of_birth"`
	Password         string    `json:"password" gorm:"column:password;size:255"`
	State            bool      `json:"state" gorm:"column:state"`
	Verify           bool      `json:"verify" gorm:"column:verify"`
	Otp              string    `json:"otp" gorm:"column:otp"`
	TimeExpiredOtp   time.Time `json:"time_expired_otp" gorm:"column:time_expired_otp"`
	Commission       float32   `json:"commission" gorm:"column:commission"`
	Deleted          bool      `json:"deleted" gorm:"column:deleted"`
	DeletedBy        int       `json:"deleted_by" gorm:"column:deleted_by"`
	UpdatedBy        int       `json:"updated_by" gorm:"column:updated_by"`
	CreatedBy        int       `json:"created_by" gorm:"column:created_by"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt        time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	InvitationStatus int       `json:"invitation_status" gorm:"-"`
	CountContract    int       `json:"count_contract" gorm:"-"`
	SumPrice         int       `json:"sum_price" gorm:"-"`
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
