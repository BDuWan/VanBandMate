package models

import "time"

type StudyProgram struct {
	StudyProgramID int       `json:"study_program_id" gorm:"primaryKey;autoIncrement"`
	Title          string    `json:"title" gorm:"size:255"`
	Image          string    `json:"image" gorm:"size:255"`
	NumberStudent  int       `json:"number_student"`
	MaxNumber      int       `json:"max_number"`
	Description    string    `json:"description" gorm:"size:255"`
	Deleted        bool      `json:"deleted"`
	DeletedBy      int       `json:"deleted_by"`
	UpdatedBy      int       `json:"updated_by"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
