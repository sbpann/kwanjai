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
		err := ginContext.ShouldBindJSON(board)
		if board.Project == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Project UUID is reqired."})
			return
		}

		// Check project membership
		getProject, err := libraries.FirestoreFind("projects", board.Project)
		if !getProject.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		project := new(models.Project)
		err = getProject.DataTo(project)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		_, found := libraries.Find(project.Members, username)
		if !found {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		searchBoards, err := libraries.FirestoreSearch("boards", "Project", "==", board.Project)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		allBoards := []*models.Board{}
		for _, board := range searchBoards {
			b := new(models.Board)
			board.DataTo(b)
			allBoards = append(allBoards, b)
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
		getProject, err := libraries.FirestoreFind("projects", board.Project)
		if !getProject.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		project := new(models.Project)
		err = getProject.DataTo(project)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if project.User != username {
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
		board := new(models.Board)
		ginContext.ShouldBindJSON(board)
		if board.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}

		// Check project membership
		username := helpers.GetUsername(ginContext)
		getProject, err := libraries.FirestoreFind("projects", board.Project)
		if !getProject.Exists() {
			ginContext.JSON(http.StatusNotFound, gin.H{"message": "Project not found."})
			return
		}
		project := new(models.Project)
		err = getProject.DataTo(project)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		_, found := libraries.Find(project.Members, username)
		if !found {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		// end

		status, message, board := board.FindBoard()
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

		// Check board ownership
		copiedBoard := new(models.Board)
		copiedBoard.UUID = board.UUID
		status, message, _ := copiedBoard.FindBoard()
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
		// end

		board.User = copiedBoard.User
		board.Project = copiedBoard.Project
		status, message, board = board.UpdateBoard("Name", board.Name)
		status, message, board = board.UpdateBoard("Description", board.Description)
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

		// Check board ownership
		copiedBoard := new(models.Board)
		copiedBoard.UUID = board.UUID
		status, message, _ := copiedBoard.FindBoard()
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
		// end

		status, message, board = board.DeleteBoard()
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
