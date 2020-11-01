package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewBoard endpoint
func NewBoard() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		board := new(models.Board)
		err := ginContext.ShouldBindJSON(board)
		board.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Board name is required."})
			return
		}
		// Check project owner
		project := new(models.Project)
		getProject, err := libraries.FirestoreFind("projects", board.Project)
		if !getProject.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		err = getProject.DataTo(project)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if project.User != username {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		status, message, board := models.NewBoard(board)
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
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		status, message, board := models.FindBoard(board)
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
		if board.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		copiedBoard := new(models.Board)
		copiedBoard.UUID = board.UUID
		status, message, _ := models.FindBoard(copiedBoard)
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedBoard.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		board.User = copiedBoard.User
		status, message, board = models.UpdateBoard(board)
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
		if board.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		copiedBoard := new(models.Board)
		copiedBoard.UUID = board.UUID
		status, message, _ := models.FindBoard(copiedBoard)
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedBoard.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		status, message, board = models.DeleteBoard(board)
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
