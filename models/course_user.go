package models

import "time"

type CourseUser struct {
	CourseUserID int       `json:"course_user_id" gorm:"primaryKey;autoIncrement"`
	CourseID     int       `json:"course_id"`
	Course       Course    `gorm:"foreignKey:CourseID;references:CourseID"`
	UserID       int       `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID;references:UserID"`
	Deleted      bool      `json:"deleted"`
	DeletedBy    int       `json:"deleted_by"`
	UpdatedBy    int       `json:"updated_by"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
