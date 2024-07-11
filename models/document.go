package models

import "time"

type Document struct {
	DocumentID  int       `json:"document_id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:255"`
	Description string    `json:"description" gorm:"size:255"`
	CourseID    int       `json:"course_id"`
	Course      Course    `gorm:"foreignKey:CourseID;references:CourseID"`
	Deleted     bool      `json:"deleted"`
	DeletedBy   int       `json:"deleted_by"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type Document1 struct {
	DocumentID  int       `json:"document_id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:255"`
	Description string    `json:"description" gorm:"size:255"`
	ClassID     int       `json:"class_id"`
	Class       Class     `gorm:"foreignKey:ClassID;references:ClassID"`
	Deleted     bool      `json:"deleted"`
	DeletedBy   int       `json:"deleted_by"`
	UpdatedBy   int       `json:"updated_by"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
