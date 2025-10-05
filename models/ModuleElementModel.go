package models

import (
	"cloud.google.com/go/datastore"
)

// JSON date string will look like this: "2021-02-18T21:54:42.123Z"
// err := time.Parse(time.RFC3339, dateString) //RFC 3339 is a profile for ISO 8601

type ModuleElement struct {
	KeyID     int64 `json:"id"` //gorm:"primary_key,autoIncrement"
	ModuleID  int64 `json:"module_id,omitempty"`
	ElementID int64 `json:"element_id,omitempty"`
	SortKey   int   `json:"sort_key,omitempty"`
}

type ModuleElementRepository interface {
	CreateModuleElement(userCourse *ModuleElement) (*datastore.Key, error)
	GetAllModuleElements() ([]*ModuleElement, error)
	DeleteModuleElement(id int64) error
	GetModuleElementsByModuleID(moduleID int64) ([]*ModuleElement, error)
	GetModuleElementsByElementID(elementID int64) ([]*ModuleElement, error)
	GetModuleElementByModuleIDAndElementID(moduleID int64, elementID int64) (*ModuleElement, error)
	UpdateModuleElement(id int64, userCourse *ModuleElement) (*datastore.Key, error)
	GetElementsByModuleID(moduleID int64) ([]*Element, error)
	GetModulesByElementID(elementID int64) ([]*Module, error)
	GetElementsByInstructorID(instructorID int64) ([]*Element, error)
}
