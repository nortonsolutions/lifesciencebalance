package models

import "cloud.google.com/go/datastore"

// create Module model
type Module struct {
	// auto increment id
	KeyID       int64   `json:"id"` //gorm:"primary_key,autoIncrement"
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty" datastore:",noindex"`
	TimeLimit   int     `json:"time_limit,omitempty"`
	MaxAttempts int     `json:"max_attempts,omitempty"`
	MinPassing  int     `json:"min_passing,omitempty"`
	SortKey     int     `json:"sort_key,omitempty"`
	CourseID    int64   `json:"course_id,omitempty"`
	ThreadIDs   []int64 `json:"thread_ids,omitempty" datastore:",noindex"`
	OwnerID     int64   `json:"owner_id,omitempty"`
}

// ModuleRepository ..
type ModuleRepository interface {
	CreateModule(Module *Module) (*datastore.Key, error)
	GetAllModules() ([]*Module, error)
	DeleteModule(id int64) error
	GetModuleByID(id int64) (*Module, error)
	UpdateModule(id int64, Module *Module) (*datastore.Key, error)
	GetAllModulesByCourseID(courseID int64) ([]*Module, error)
}
