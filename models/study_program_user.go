package models

import "time"

type StudyProgramUser struct {
	StudyProgramUserID int          `json:"study_program_user_id" gorm:"primaryKey;autoIncrement"`
	StudyProgramID     int          `json:"study_program_id"`
	StudyProgram       StudyProgram `gorm:"foreignKey:StudyProgramID;references:StudyProgramID"`
	UserID             int          `json:"user_id"`
	User               User         `gorm:"foreignKey:UserID;references:UserID"`
	Deleted            bool         `json:"deleted"`
	DeletedBy          int          `json:"deleted_by"`
	UpdatedBy          int          `json:"updated_by"`
	CreatedBy          int          `json:"created_by"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	DeletedAt          time.Time    `json:"deleted_at"`
}
