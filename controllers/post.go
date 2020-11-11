package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AllPost endpoint
func AllPost() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.Project == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Project UUID is reqired."})
			return
		}

		searchPosts, err := libraries.FirestoreSearch("posts", "Project", "==", post.Project)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		allPosts := []*models.Post{}
		for _, p := range searchPosts {
			post := new(models.Post)
			p.DataTo(post)
			post.ID = p.Ref.ID
			// check project membership
			if !helpers.IsProjectMember(username, post.Project) {
				ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
				return
			}
			allPosts = append(allPosts, post)
		}
		ginContext.JSON(http.StatusOK,
			gin.H{
				"posts": allPosts,
			})
	}
}

// NewPost endpoint
func NewPost() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		err := ginContext.ShouldBindJSON(post)
		post.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		getBoard, err := libraries.FirestoreFind("boards", post.Board)
		if err != nil {
			log.Panicln(err)
		}
		if !getBoard.Exists() {
			ginContext.JSON(http.StatusNotFound, "Board not found.")
			return
		}
		board := new(models.Board)
		getBoard.DataTo(board)

		// Check project membership
		post.Project = board.Project
		if !helpers.IsProjectMember(username, post.Project) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		post.Project = board.Project
		status, message, post := post.CreatePost()
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"post":    post,
			})
	}
}

// FindPost endpoint
func FindPost() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		status, message, post := post.FindPost()

		// Check project membership
		if !helpers.IsProjectMember(username, post.Project) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		ginContext.JSON(status,
			gin.H{
				"message": message,
				"post":    post,
			})
	}
}

// UpdatePost endpoint
func UpdatePost() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID."})
			return
		}

		// Check post ownership
		if !helpers.IsOwner(username, "posts", post.ID) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		_, err := libraries.FirestoreFind("boards", post.Board)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid board ID."})
			return
		}
		status, message, _ := post.UpdatePost("Board", post.Board)
		status, message, _ = post.UpdatePost("Title", post.Title)
		status, message, _ = post.UpdatePost("Content", post.Content)
		status, message, _ = post.UpdatePost("Completed", post.Completed)
		status, message, _ = post.UpdatePost("Urgent", post.Urgent)
		status, message, _ = post.UpdatePost("People", post.People)
		status, message, _ = post.UpdatePost("LastModified", time.Now().Truncate(time.Millisecond))
		status, message, _ = post.UpdatePost("DueDate", post.DueDate)
		status, message, _ = post.FindPost()
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"post":    post,
			})
	}
}

// DeletePost endpoint
func DeletePost() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		post := new(models.Post)
		ginContext.ShouldBindJSON(post)
		if post.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check post ownership
		if !helpers.IsOwner(username, "posts", post.ID) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		status, message, _ := post.DeletePost()
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
