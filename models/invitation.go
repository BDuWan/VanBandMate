package models

import "time"

type Invitation struct {
	InvitationID int        `gorm:"primaryKey;autoIncrement;column:invitation_id" json:"invitation_id"`
	HiringNewsID int        `gorm:"column:hiring_news_id" json:"hiring_news_id"`
	HiringNews   HiringNews `gorm:"foreignKey:HiringNewsID;references:HiringNewsID"`
	NhaccongID   int        `gorm:"column:nhaccong_id" json:"nhaccong_id"`
	Status       int        `gorm:"column:status" json:"status"`
	InviteAt     time.Time  `gorm:"column:invite_at" json:"invite_at"`
}

//Status
//0: đã gửi, đợi phản hồi
//1: đã chấp nhận
//2: đã hủy (chủ loa đài hủy hoặc hủy tự động)
//3: bị từ chối (nhạc công từ chối)
//4: chưa gửi
//5: trùng lịch
