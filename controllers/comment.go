package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NewComment endpoint
func NewComment() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.UUID == "" || len(post.Comments) != 1 || post.Comments[0].Body == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid comment form."})
			return
		}
		// Copy data before find post
		comment := new(models.Comment)
		comment.Body = post.Comments[0].Body
		comment.UUID = uuid.New().String()
		comment.User = username
		now := time.Now().Truncate(time.Millisecond)
		comment.AddedDate = now
		comment.LastModified = now
		status, message, post := post.FindPost()
		if status != http.StatusOK {
			ginContext.JSON(status, gin.H{"message": message})
			return
		}
		// Append new comment to old post
		post.Comments = append(post.Comments, comment)

		// Check project membership
		if !helpers.IsProjectMember(username, post.Project) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		// Update to Firestore
		_, err := libraries.FirestoreUpdateField("posts", post.UUID, "Comments", post.Comments)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ginContext.JSON(http.StatusCreated,
			gin.H{
				"message": "Created sucessfully",
				"post":    post,
			})
	}
}

// UpdateComment endpoint
func UpdateComment() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.UUID == "" ||
			len(post.Comments) != 1 ||
			post.Comments[0].Body == "" ||
			post.Comments[0].UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid comment form."})
			return
		}
		// Copy data before find post
		comment := new(models.Comment)
		comment.UUID = post.Comments[0].UUID
		comment.Body = post.Comments[0].Body
		comment.LastModified = time.Now().Truncate(time.Millisecond)
		status, message, post := post.FindPost()
		if status != http.StatusOK {
			ginContext.JSON(status, gin.H{"message": message})
			return
		}

		// Edit comment
		var owner string
		var found bool
		for _, c := range post.Comments {
			if c.UUID == comment.UUID {
				c.Body = comment.Body
				c.LastModified = comment.LastModified
				owner = c.User
				found = true
				break
			}
			found = false
		}
		if !found {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
			return
		}
		if username != owner {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		// Update to Firestore
		_, err := libraries.FirestoreUpdateField("posts", post.UUID, "Comments", post.Comments)
		if err != nil {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		ginContext.JSON(status,
			gin.H{
				"message": "Updated sucessfully.",
				"post":    post,
			})
	}
}

// DeleteComment endpoint
func DeleteComment() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.UUID == "" ||
			len(post.Comments) != 1 ||
			post.Comments[0].UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid comment form."})
			return
		}
		// Copy data before find post
		comment := new(models.Comment)
		comment.UUID = post.Comments[0].UUID
		status, message, post := post.FindPost()
		if status != http.StatusOK {
			ginContext.JSON(status, gin.H{"message": message})
			return
		}

		// Edit comment
		var owner string
		var found bool
		var index int
		for i, c := range post.Comments {
			if c.UUID == comment.UUID {
				c.Body = comment.Body
				c.LastModified = comment.LastModified
				owner = c.User
				found = true
				index = i
				break
			}
			found = false
		}
		if !found {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
			return
		}
		if username != owner {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		post.Comments = append(post.Comments[:index], post.Comments[index+1:]...)
		// Update to Firestore
		_, err := libraries.FirestoreUpdateField("posts", post.UUID, "Comments", post.Comments)
		if err != nil {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		ginContext.JSON(status,
			gin.H{
				"message": "Deleted sucessfully.",
				"post":    post,
			})
	}
}
