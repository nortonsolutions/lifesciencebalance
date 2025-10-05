package models

import "cloud.google.com/go/datastore"

type Course struct {
	KeyID       int64   `json:"id"` //gorm:"primary_key,autoIncrement"
	Name        string  `json:"name,omitempty"`
	HomeContent string  `json:"home_content,omitempty" datastore:",noindex"`
	Description string  `json:"description,omitempty" datastore:",noindex"`
	Modules     []int64 `json:"modules,omitempty" datastore:",noindex"`
	OwnerID     int64   `json:"owner_id,omitempty"`
	Approved    bool    `json:"approved"`
	Department  string  `json:"department,omitempty"`
}

type CourseRepository interface {
	CreateCourse(Course *Course) (*datastore.Key, error)
	GetAllCourses() ([]*Course, error)
	DeleteCourse(id int64) error
	GetCourseByID(id int64) (*Course, error)
	UpdateCourse(id int64, Course *Course) (*datastore.Key, error)
	GetApprovedCourses() ([]*Course, error)
	GetUnapprovedCourses() ([]*Course, error)
	GetCoursesByDepartment(department string) ([]*Course, error)
}
