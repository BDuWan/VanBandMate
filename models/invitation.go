package models

import "time"

type Invitation struct {
	InvitationID int       `gorm:"primaryKey;autoIncrement;column:invitation_id" json:"invitation_id"`
	HiringNewsID int       `gorm:"not null;default:0;column:hiring_news_id" json:"hiring_news_id"`
	NhaccongID   int       `gorm:"not null;default:0;column:nhaccong_id" json:"nhaccong_id"`
	Status       int       `gorm:"column:status" json:"status"`
	InviteAt     time.Time `gorm:"column:invite_at" json:"invite_at"`
}
