package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Project struct {
	KeyID       int64     `json:"id"` //gorm:"primary_key,autoIncrement"
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	File        string    `json:"file,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	CourseID    int64     `json:"course_id,omitempty"`
	ModuleID    int64     `json:"module_id,omitempty"`
}

type ProjectRepository interface {
	CreateProject(Project *Project) (*datastore.Key, error)
	GetAllProjects() ([]*Project, error)
	DeleteProject(id int64) error
	GetProjectByID(id int64) (*Project, error)
	UpdateProject(id int64, Project *Project) (*datastore.Key, error)
	GetProjectsByCourseID(courseID int64) ([]*Project, error)
	GetProjectsByModuleID(moduleID int64) ([]*Project, error)
	GetProjectsByUserID(userID int64) ([]*Project, error)
	GetProjectByUserIDandModuleID(userID int64, moduleID int64) (*Project, error)
}
