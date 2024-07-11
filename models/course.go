package models

import "time"

// type StudyProgram struct {
// 	StudyProgramID    int       `json:"course_id" gorm:"primaryKey;autoIncrement"`
// 	Title       string    `json:"title" gorm:"size:255"`
// 	Image       string    `json:"image" gorm:"size:255"`
// 	Price       string    `json:"price" gorm:"size:100"`
// 	Description string    `json:"description" gorm:"size:255"`
// 	Deleted     bool      `json:"deleted"`
// 	DeletedBy   int       `json:"deleted_by"`
// 	UpdatedBy   int       `json:"updated_by"`
// 	CreatedBy   int       `json:"created_by"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	DeletedAt   time.Time `json:"deleted_at"`
// }

type Course struct {
	CourseID       int          `json:"course_id" gorm:"primaryKey;autoIncrement"`
	Name           string       `json:"name" gorm:"size:255"`
	UserID         int          `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;references:UserID"`
	StudyProgramID int          `json:"study_program_id"`
	StudyProgram   StudyProgram `gorm:"foreignKey:StudyProgramID;references:StudyProgramID"`
	Description    string       `json:"description" gorm:"size:255"`
	Image          string       `json:"image" gorm:"size:255"`
	NumberStudent  int          `json:"number_student"`
	Deleted        bool         `json:"deleted"`
	DeletedBy      int          `json:"deleted_by"`
	UpdatedBy      int          `json:"updated_by"`
	CreatedBy      int          `json:"created_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      time.Time    `json:"deleted_at"`
}
