package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// anything that could possibly be the answer to a question
type Answer struct {
	Answer      []bool `json:"answer,omitempty" datastore:",noindex"`
	AnswerText  string `json:"answer_text,omitempty"`
	AnswerEssay string `json:"answer_essay,omitempty" datastore:",noindex"`
	ProjectID   int64  `json:"project_id,omitempty"`
	Correct     bool   `json:"correct,omitempty"`
}

// the string key for Answers is the element ID
type UserModule struct {
	UserID     int64             `json:"user_id,omitempty"`
	ModuleID   int64             `json:"module_id,omitempty"`
	Answers    map[string]Answer `json:"answers,omitempty" datastore:",noindex"`
	Date       string            `json:"date,omitempty"`
	Score      int               `json:"score,omitempty"`
	TimePassed int               `json:"time_passed,omitempty"`
}

// create User model
type User struct {
	// auto increment id
	KeyID     int64        `json:"id"` //gorm:"primary_key,autoIncrement"
	Username  string       `json:"username,omitempty"`
	Email     string       `json:"email,omitempty"`
	Password  string       `json:"password,omitempty"`
	Firstname string       `json:"firstname,omitempty"`
	Lastname  string       `json:"lastname,omitempty"`
	Roles     []string     `json:"roles,omitempty" datastore:",noindex"`
	Modules   []UserModule `json:"modules,omitempty" datastore:",noindex"`
	Bio       string       `json:"bio,omitempty"`
	Avatar    string       `json:"avatar,omitempty"`
	CreatedOn time.Time    `json:"created_on,omitempty"`
}

// TODO: why is this in the model?
func (u User) GetRoles() []string {
	if u.Roles == nil {
		return []string{}
	}
	return u.Roles
}

// UserRepository ..
type UserRepository interface {
	CreateUser(user *User) (*datastore.Key, error)
	GetAllUsers() ([]*User, error)
	GetUserByID(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(id int64, user *User) (*datastore.Key, error)
	DeleteUser(id int64) error
	GetUserByUsernameAndPassword(username string, password string) (*User, error)
}
