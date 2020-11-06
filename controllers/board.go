package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AllBoard endpoint
func AllBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.Project == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Project UUID is reqired."})
			return
		}

		searchBoards, err := libraries.FirestoreSearch("boards", "Project", "==", board.Project)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		allBoards := []*models.Board{}
		for _, b := range searchBoards {
			board := new(models.Board)
			b.DataTo(board)
			// check project membership
			if !helpers.IsProjectMember(username, board.Project) {
				ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
				return
			}
			allBoards = append(allBoards, board)
		}
		ginContext.JSON(http.StatusOK,
			gin.H{
				"boards": allBoards,
			})
	}
}

// NewBoard endpoint
func NewBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		err := ginContext.ShouldBindJSON(board)
		board.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Check project ownership
		if !helpers.IsOwner(username, "projects", board.Project) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		status, message, board := board.CreateBoard()
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"board":   board,
			})
	}
}

// FindBoard endpoint
func FindBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID."})
			return
		}
		status, message, board := board.FindBoard()

		// Check project membership
		if !helpers.IsProjectMember(username, board.Project) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		ginContext.JSON(status,
			gin.H{
				"message": message,
				"board":   board,
			})
	}
}

// UpdateBoard endpoint
func UpdateBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check board ownership
		if !helpers.IsOwner(username, "boards", board.ID) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		status, message, board := board.UpdateBoard("Name", board.Name)
		status, message, board = board.UpdateBoard("Description", board.Description)
		status, message, board = board.FindBoard()
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"board":   board,
			})
	}
}

// DeleteBoard endpoint
func DeleteBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.ID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check board ownership
		if !helpers.IsOwner(username, "boards", board.ID) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}

		status, message, _ := board.DeleteBoard()
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
