package models

import (
	"time"
)

type LessonUser struct {
	LessonUserID int       `json:"lesson_user_id" gorm:"primaryKey;autoIncrement"`
	LessonID     int       `json:"lesson_id"`
	Lesson       Lesson    `gorm:"foreignKey:LessonID;references:LessonID"`
	UserID       int       `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID;references:UserID"`
	Result       string    `json:"result" gorm:"size:255"`
	Status       int       `json:"status"`
	Deleted      bool      `json:"deleted"`
	DeletedBy    int       `json:"deleted_by"`
	UpdatedBy    int       `json:"updated_by"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

