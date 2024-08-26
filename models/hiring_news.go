package models

import (
	"time"
)

type HiringNews struct {
	HiringNewsID    int              `json:"hiring_news_id" gorm:"primaryKey;autoIncrement"`
	ChuloadaiID     int              `json:"chuloadai_id"`
	User            User             `gorm:"foreignKey:ChuloadaiID;references:UserID"`
	ProvinceCode    string           `json:"province_code"`
	Province        Province         `gorm:"foreignKey:ProvinceCode;references:Code"`
	DistrictCode    string           `json:"district_code"`
	District        District         `gorm:"foreignKey:DistrictCode;references:Code"`
	WardCode        string           `json:"ward_code"`
	Ward            Ward             `gorm:"foreignKey:WardCode;references:Code"`
	AddressDetail   string           `json:"address_detail"`
	Date            time.Time        `json:"date"`
	Describe        string           `json:"describe"`
	Price           int              `json:"price"`
	HiringEnough    bool             `json:"hiring_enough"`
	Deleted         bool             `json:"deleted"`
	DeletedBy       int              `json:"deleted_by"`
	UpdatedBy       int              `json:"updated_by"`
	CreatedBy       int              `json:"created_by"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       time.Time        `json:"deleted_at"`
	Applicants      []UserHiringNews `json:"applicants" gorm:"foreignKey:HiringNewsID;references:HiringNewsID"`
	ApplicantStatus int              `json:"applicant_status" gorm:"-"`
}

//ApplicantStatus:
//0: đang ứng tuyển nhưng chưa dc chấp nhận
//1: đã được chấp nhận
//2: đã ứng tuyển nhưng đã hủy (cá nhân chủ động hủy)
//3: đã ứng tuyển nhưng đã hủy (tự động hủy do cơ chế tạo hợp đồng)
//4: chưa ứng tuyển
//5: trùng lịch
