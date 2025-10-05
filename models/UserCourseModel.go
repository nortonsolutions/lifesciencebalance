package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// JSON date string will look like this: "2021-02-18T21:54:42.123Z"
// err := time.Parse(time.RFC3339, dateString) //RFC 3339 is a profile for ISO 8601

type UserCourse struct {
	KeyID       int64     `json:"id"` //gorm:"primary_key,autoIncrement"
	UserID      int64     `json:"user_id,omitempty"`
	CourseID    int64     `json:"course_id,omitempty"`
	Grade       int       `json:"grade,omitempty"`
	StartedOn   time.Time `json:"started_on,omitempty"`
	CompletedOn time.Time `json:"completed_on,omitempty"`
	Role        string    `json:"role,omitempty"`
}

type UserCourseRepository interface {
	CreateUserCourse(userCourse *UserCourse) (*datastore.Key, error)
	GetAllUserCourses() ([]*UserCourse, error)
	DeleteUserCourse(id int64) error
	GetUserCoursesByUserID(userID int64) ([]*UserCourse, error)
	GetUserCoursesByCourseID(courseID int64) ([]*UserCourse, error)
	GetUserCourseByUserIDAndCourseID(userID int64, courseID int64) (*UserCourse, error)
	UpdateUserCourse(id int64, userCourse *UserCourse) (*datastore.Key, error)
	GetCoursesByUserID(userID int64) ([]*Course, error)
	GetUsersByCourseID(courseID int64) ([]*User, error)
	GetInstructorsByCourseID(courseID int64) ([]*User, error)
}
