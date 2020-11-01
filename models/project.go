package models

import (
	"kwanjai/libraries"
	"net/http"

	"github.com/google/uuid"
)

// Project model
type Project struct {
	UUID string `json:"uuid"`
	User string `json:"username"`
	Name string `json:"name" binding:"required"`
}

type projectPerform interface {
	createProject() (int, string, *Project)
	findProject() (int, string, *Project)
	updateProject() (int, string, *Project)
	deleteProject() (int, string, *Project)
}

// NewProject project method for interface with controller.
func NewProject(perform projectPerform) (int, string, *Project) {
	status, message, project := perform.createProject()
	return status, message, project
}

// FindProject project method for interface with controller.
func FindProject(perform projectPerform) (int, string, *Project) {
	status, message, project := perform.findProject()
	return status, message, project
}

// UpdateProject project method for interface with controller.
func UpdateProject(perform projectPerform) (int, string, *Project) {
	status, message, project := perform.updateProject()
	return status, message, project
}

// DeleteProject project method for interface with controller.
func DeleteProject(perform projectPerform) (int, string, *Project) {
	status, message, project := perform.deleteProject()
	return status, message, project
}

func (project *Project) createProject() (int, string, *Project) {
	project.UUID = uuid.New().String()
	_, err := libraries.FirestoreCreatedOrSet("projects", project.UUID, project)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusCreated, "Created project.", project
}

func (project *Project) findProject() (int, string, *Project) {
	if project.UUID == "" {
		return http.StatusNotFound, "Project not found.", nil
	}
	getProject, _ := libraries.FirestoreFind("projects", project.UUID)
	if getProject.Exists() {
		getProject.DataTo(project)
		return http.StatusOK, "Get project successfully.", project
	}
	return http.StatusNotFound, "Project not found.", nil
}

func (project *Project) updateProject() (int, string, *Project) {
	_, err := libraries.FirestoreUpdateField("projects", project.UUID, "Name", project.Name)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", project
}

func (project *Project) deleteProject() (int, string, *Project) {
	_, err := libraries.FirestoreDelete("projects", project.UUID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
