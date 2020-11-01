package controllers

import (
	"kwanjai/helpers"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewProject endpoint
func NewProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		project := new(models.Project)
		err := ginContext.ShouldBindJSON(project)
		project.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Project name is required."})
			return
		}
		status, message, project := models.NewProject(project)
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"project": project,
			})
	}
}

// FindProject endpoint
func FindProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		project := new(models.Project)
		ginContext.ShouldBindJSON(project)
		if project.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		status, message, project := models.FindProject(project)
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"project": project,
			})
	}
}

// UpdateProject endpoint
func UpdateProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		project := new(models.Project)
		ginContext.ShouldBindJSON(project)
		if project.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		copiedProject := new(models.Project)
		copiedProject.UUID = project.UUID
		status, message, _ := models.FindProject(copiedProject)
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedProject.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		project.User = copiedProject.User
		status, message, project = models.UpdateProject(project)
		ginContext.JSON(status,
			gin.H{
				"message": message,
				"project": project,
			})
	}
}

// DeleteProject endpoint
func DeleteProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		project := new(models.Project)
		ginContext.ShouldBindJSON(project)
		if project.UUID == "" {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID."})
			return
		}
		copiedProject := new(models.Project)
		copiedProject.UUID = project.UUID
		status, message, _ := models.FindProject(copiedProject)
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedProject.User != username || helpers.IsSuperUser(ginContext) {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		status, message, project = models.DeleteProject(project)
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
