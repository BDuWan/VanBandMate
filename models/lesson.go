package models

import "time"

type Lesson struct {
	LessonID    int    `json:"lesson_id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"size:255"`
	Description string `json:"description" gorm:"size:255"`
	CourseID    int    `json:"course_id"`
	Course      Course `gorm:"foreignKey:CourseID;references:CourseID"`
	//StartTime   string    `json:"start_time" gorm:"size:50"`
	//EndTime     string    `json:"end_time" gorm:"size:50"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	LinkRecord string    `json:"link_record" gorm:"size:255"`
	LinkStudy  string    `json:"link_study" gorm:"size:255"`
	Deleted    bool      `json:"deleted"`
	DeletedBy  int       `json:"deleted_by"`
	UpdatedBy  int       `json:"updated_by"`
	CreatedBy  int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type Lesson1 struct {
	LessonID    int    `json:"lesson_id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"size:255"`
	Description string `json:"description" gorm:"size:255"`
	ClassID     int    `json:"class_id"`
	Class       Class  `gorm:"foreignKey:ClassID;references:ClassID"`
	//StartTime   string    `json:"start_time" gorm:"size:50"`
	//EndTime     string    `json:"end_time" gorm:"size:50"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	LinkRecord string    `json:"link_record" gorm:"size:255"`
	LinkStudy  string    `json:"link_study" gorm:"size:255"`
	Deleted    bool      `json:"deleted"`
	DeletedBy  int       `json:"deleted_by"`
	UpdatedBy  int       `json:"updated_by"`
	CreatedBy  int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
