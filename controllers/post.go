package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
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
			p.DataTo(post)
			// check project membership
			if helpers.IsProjectMember(username, post.Project) {
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
			ginContext.JSON(http.StatusInternalServerError, err.Error())
		}
		if !getBoard.Exists() {
			ginContext.JSON(http.StatusNotFound, "Board not found.")
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
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check post ownership
		if !helpers.IsOwner(username, "post", post.ID) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		status, message, post := post.UpdatePost("Title", post.Title)
		status, message, post = post.UpdatePost("Body", post.Body)
		status, message, post = post.UpdatePost("Completed", post.Completed)
		status, message, post = post.UpdatePost("Urgent", post.Urgent)
		status, message, post = post.UpdatePost("People", post.People)
		status, message, post = post.UpdatePost("LastModified", time.Now().Truncate(time.Millisecond))
		status, message, post = post.FindPost()
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
