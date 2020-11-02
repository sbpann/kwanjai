package controllers

import (
	"kwanjai/helpers"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AllProject endpoint
func AllProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		searchProjects, err := libraries.FirestoreSearch("projects", "Members", "array-contains", username)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		project := new(models.Project)
		allProjects := []*models.Project{}
		for _, p := range searchProjects {
			p.DataTo(project)
			allProjects = append(allProjects, project)
		}
		ginContext.JSON(http.StatusOK, gin.H{"projects": allProjects})
		return
	}
}

// NewProject endpoint
func NewProject() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		username := helpers.GetUsername(ginContext)
		project := new(models.Project)
		err := ginContext.ShouldBindJSON(project)
		project.User = username
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		status, message, project := project.CreateProject()
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
		status, message, project := project.FindProject()
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
		status, message, _ := copiedProject.FindProject()
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedProject.User != username {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		project.User = copiedProject.User
		status, message, project = project.UpdateProject()
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
		status, message, _ := copiedProject.FindProject()
		if status != http.StatusOK {
			ginContext.JSON(status,
				gin.H{
					"message": message,
				})
			return
		}
		if copiedProject.User != username {
			ginContext.JSON(http.StatusForbidden, gin.H{"message": "You cannot perform this action."})
			return
		}
		status, message, project = project.DeleteProject()
		ginContext.JSON(status,
			gin.H{
				"message": message,
			})
	}
}
