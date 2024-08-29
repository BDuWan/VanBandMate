package models

import "time"

type UserHiringNews struct {
	UserHiringNewsID int        `json:"user_hiring_news_id" gorm:"primaryKey;autoIncrement"`
	HiringNewsID     int        `json:"hiring_news_id"`
	HiringNews       HiringNews `gorm:"foreignKey:HiringNewsID;references:HiringNewsID"`
	NhaccongID       int        `json:"nhaccong_id"`
	User             User       `gorm:"foreignKey:NhaccongID;references:UserID"`
	Status           int        `json:"status"`
	ApplyAt          time.Time  `json:"apply_at"`
	Date             time.Time  `json:"date"`
}

//Status
//0: Đang ứng tuyển, chờ chấp nhận
//1: Đã chấp nhận.
//2: Đã từng ứng tuyển nhưng thu hồi (nhạc công chủ động)
//3: Đã từng ứng tuyển nhưng bị hủy tự động (do cơ chế tạo hợp đồng)
