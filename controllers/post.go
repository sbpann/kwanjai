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
		err := ginContext.ShouldBindJSON(post)
		post.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Check project membership
		getBoard, err := libraries.FirestoreFind("boards", post.Board)
		if !getBoard.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		board := new(models.Board)
		getBoard.DataTo(board)
		project := new(models.Project)
		getProject, _ := libraries.FirestoreFind("projects", board.Project) // board existence ensures no error.
		getProject.DataTo(project)
		_, found := libraries.Find(project.Members, username)
		if !found {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		searchPosts, err := libraries.FirestoreSearch("post", "Board", "==", post.Board)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		allPosts := []*models.Post{}
		for _, post := range searchPosts {
			p := new(models.Post)
			post.DataTo(p)
			allPosts = append(allPosts, p)
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

		// Check project membership
		getBoard, err := libraries.FirestoreFind("boards", post.Board)
		if !getBoard.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		board := new(models.Board)
		getBoard.DataTo(board)
		project := new(models.Project)
		getProject, _ := libraries.FirestoreFind("projects", board.Project) // board existence ensures no error.
		getProject.DataTo(project)
		_, found := libraries.Find(project.Members, username)
		if !found {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

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
		if post.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check project membership
		getBoard, err := libraries.FirestoreFind("boards", post.Board)
		if !getBoard.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		board := new(models.Board)
		getBoard.DataTo(board)
		project := new(models.Project)
		getProject, _ := libraries.FirestoreFind("projects", board.Project) // board existence ensures no error.
		getProject.DataTo(project)
		_, found := libraries.Find(project.Members, username)
		if !found {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		status, message, post := post.FindPost()
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
		if post.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check post ownership
		copiedPost := new(models.Post)
		copiedPost.UUID = post.UUID
		status, message, _ := copiedPost.FindPost()
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedPost.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		post.User = copiedPost.User
		post.Board = copiedPost.Board
		status, message, post = post.UpdatePost("Title", post.Title)
		status, message, post = post.UpdatePost("Body", post.Body)
		status, message, post = post.UpdatePost("Completed", post.Completed)
		status, message, post = post.UpdatePost("Urgent", post.Urgent)
		status, message, post = post.UpdatePost("Body", post.Body)
		status, message, post = post.UpdatePost("LastModified", time.Now().Truncate(time.Millisecond))
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
		if post.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check post ownership
		copiedPost := new(models.Post)
		copiedPost.UUID = post.UUID
		status, message, _ := copiedPost.FindPost()
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedPost.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		status, message, _ = post.DeletePost()
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
