package models

import "time"

type Comment struct {
	CommentID     int       `json:"comment_id" gorm:"primaryKey"`
	UserID        int       `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID;references:UserID"`
	ParentCmtID   int       `json:"parent_cmt_id"`
	ParentComment *Comment  `gorm:"foreignKey:ParentCmtID;references:CommentID"`
	CourseID      int       `json:"course_id"`
	Course        Course    `gorm:"foreignKey:CourseID;references:CourseID"`
	IsChildCmt    bool      `json:"is_child_cmt"`
	UserIDReply   int       `json:"user_id_reply"`
	Content       string    `json:"content"`
	Deleted       bool      `json:"deleted"`
	UpdatedBy     int       `json:"updated_by"`
	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
	SubComments   []Comment `json:"sub_comments" gorm:"-"`
}
