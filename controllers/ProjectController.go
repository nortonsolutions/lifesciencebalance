package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ProjectHandler ..
type ProjectHandler struct {
	projectRepository models.ProjectRepository
}

// NewProjectHandler ..
func NewProjectHandler(projectRepository models.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepository: projectRepository}
}

// add project
func (c *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	project := models.Project{}
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.projectRepository.CreateProject(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// remove project
func (c *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.projectRepository.DeleteProject(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted"})
}

// get all projects
func (c *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := c.projectRepository.GetAllProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// get project by id
func (c *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	project, err := c.projectRepository.GetProjectByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// update project
func (c *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	project := models.Project{}
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.projectRepository.UpdateProject(idInt, &project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

func (c *ProjectHandler) GetProjectsByCourseID(w http.ResponseWriter, r *http.Request) {
	var courseID = mux.Vars(r)["courseID"]
	courseIDInt, _ := strconv.ParseInt(courseID, 10, 64)
	projects, err := c.projectRepository.GetProjectsByCourseID(courseIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (c *ProjectHandler) GetProjectsByUserID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	projects, err := c.projectRepository.GetProjectsByUserID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (c *ProjectHandler) GetProjectByUserIDandModuleID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["userID"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	var moduleID = mux.Vars(r)["moduleID"]
	moduleIDInt, _ := strconv.ParseInt(moduleID, 10, 64)
	project, err := c.projectRepository.GetProjectByUserIDandModuleID(userIDInt, moduleIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (c *ProjectHandler) GetProjectsByModuleID(w http.ResponseWriter, r *http.Request) {
	var moduleID = mux.Vars(r)["moduleID"]
	moduleIDInt, _ := strconv.ParseInt(moduleID, 10, 64)
	projects, err := c.projectRepository.GetProjectsByModuleID(moduleIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
