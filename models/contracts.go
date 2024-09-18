package models

import (
	"time"
)

// Contract represents the structure of the 'contracts' table
type Contract struct {
	ContractID      int       `gorm:"primaryKey;column:contract_id" json:"contract_id"`
	ChuloadaiID     int       `gorm:"column:chuloadai_id" json:"chuloadai_id"`
	ChuLoaDai       User      `gorm:"foreignKey:ChuloadaiID;references:UserID"`
	NhaccongID      int       `gorm:"column:nhaccong_id" json:"nhaccong_id"`
	NhacCong        User      `gorm:"foreignKey:NhaccongID;references:UserID"`
	ProvinceCode    string    `gorm:"column:province_code" json:"province_code"`
	Province        Province  `json:"province" gorm:"foreignKey:ProvinceCode;references:Code"`
	DistrictCode    string    `gorm:"column:district_code" json:"district_code"`
	District        District  `json:"district" gorm:"foreignKey:DistrictCode;references:Code"`
	WardCode        string    `gorm:"column:ward_code" json:"ward_code"`
	Ward            Ward      `json:"ward" gorm:"foreignKey:WardCode;references:Code"`
	AddressDetail   string    `gorm:"column:address_detail" json:"address_detail"`
	Price           int       `gorm:"column:price" json:"price"`
	Date            time.Time `gorm:"column:date" json:"date"`
	Status          int       `gorm:"column:status" json:"status"`
	RequestDeleteBy int       `gorm:"column:request_delete_by" json:"request_delete_by"`
	Deleted         bool      `gorm:"column:deleted" json:"deleted"`
	DeletedBy       int       `gorm:"column:deleted_by" json:"deleted_by"`
	UpdatedBy       int       `gorm:"column:updated_by" json:"updated_by"`
	CreatedBy       int       `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

//Status
//0: đã hoàn thành
//1: chưa hoàn thành
//2: có yêu cầu hủy
//3: đã hủy
