package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zetamatta/go-outputdebug"
	"lms/initializers"
	"lms/models"
	"time"
)

func APIGetComment(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIGetComment")
	courseId := c.Params("id")
	DB := initializers.DB
	var comments []models.Comment

	if err := DB.Model(&models.Comment{}).Joins("User").Where(
		"comments.deleted", false).Where(
		"User.deleted", false).Where(
		"comments.course_id", courseId).Find(&comments).Error; err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: " + err.Error())
	}

	commentMap := make(map[int][]models.Comment)
	parentComments := make(map[int]models.Comment)

	for _, comment := range comments {
		if comment.IsChildCmt {
			commentMap[comment.ParentCmtID] = append(commentMap[comment.ParentCmtID], comment)
		} else {
			parentComments[comment.CommentID] = comment
		}
	}

	var mapComments []models.Comment
	for parentID, childComments := range commentMap {
		if parentComment, exists := parentComments[parentID]; exists {
			parentComment.SubComments = childComments
			mapComments = append(mapComments, parentComment)
		}
	}

	for _, parentComment := range parentComments {
		if _, exists := commentMap[parentComment.CommentID]; !exists {
			mapComments = append(mapComments, parentComment)
		}
	}

	return c.JSON(mapComments)
}

func APIAddComment(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIAddComment")
	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Cannot parser JSON")
	}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := initializers.DB.Create(&comment).Error; err != nil {
		return c.JSON("Cannot create comment")
	}

	return c.JSON("Success")
}

func APIAddCommentReply(c *fiber.Ctx) error {
	outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: APIAddCommentReply")
	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		outputdebug.String(time.Now().Format("02-01-2006 15:04:05") + " [LMS]: Cannot parser JSON")
	}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := initializers.DB.Create(&comment).Error; err != nil {
		return c.JSON("Cannot create comment")
	}

	return c.JSON("Success")
}
