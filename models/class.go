package models

import "time"

// type Class1 struct {
// 	ClassID     int       `json:"class_id" gorm:"primaryKey;autoIncrement"`
// 	Name        string    `json:"name" gorm:"size:255"`
// 	UserID      int       `json:"user_id"`
// 	User        User      `gorm:"foreignKey:UserID;references:UserID"`
// 	CourseID    int       `json:"course_id"`
// 	Course      Course    `gorm:"foreignKey:CourseID;references:CourseID"`
// 	Description string    `json:"description" gorm:"size:255"`
// 	Deleted     bool      `json:"deleted"`
// 	DeletedBy   int       `json:"deleted_by"`
// 	UpdatedBy   int       `json:"updated_by"`
// 	CreatedBy   int       `json:"created_by"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	DeletedAt   time.Time `json:"deleted_at"`
// }

type Class struct {
	ClassID        int          `json:"class_id" gorm:"primaryKey;autoIncrement"`
	Name           string       `json:"name" gorm:"size:255"`
	Price          string       `json:"price" gorm:"size:100"`
	UserID         int          `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;references:UserID"`
	StudyProgramID int          `json:"study_program_id"`
	StudyProgram   StudyProgram `gorm:"foreignKey:StudyProgramID;references:StudyProgramID"`
	Description    string       `json:"description" gorm:"size:255"`
	NumberStudent  int          `json:"number_student"`
	Deleted        bool         `json:"deleted"`
	DeletedBy      int          `json:"deleted_by"`
	UpdatedBy      int          `json:"updated_by"`
	CreatedBy      int          `json:"created_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      time.Time    `json:"deleted_at"`
}
