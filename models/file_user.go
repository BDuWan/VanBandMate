package models

import "time"

type FileUser struct {
	FileUserID   int       `json:"file_user_id" gorm:"primaryKey;autoIncrement"`
	AssignmentID int       `json:"assignment_id"`
	UserID       int       `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID;references:UserID"`
	File         string    `json:"file" gorm:"size:255"`
	Deleted      bool      `json:"deleted"`
	DeletedBy    int       `json:"deleted_by"`
	UpdatedBy    int       `json:"updated_by"`
	CreatedBy    int       `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
