package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Reply struct {
	Body      string    `json:"body,omitempty"`
	Author    string    `json:"author,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	BumpedOn  time.Time `json:"bumped_on,omitempty"`
	Reported  bool      `json:"reported,omitempty"`
	Upvotes   int       `json:"upvotes,omitempty"`
	Downvotes int       `json:"downvotes,omitempty"`
}

type Thread struct {
	KeyID     int64     `json:"id"` //gorm:"primary_key,autoIncrement"
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Author    string    `json:"author,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	BumpedOn  time.Time `json:"bumped_on,omitempty"`
	Closed    bool      `json:"closed,omitempty"`
	Reported  bool      `json:"reported,omitempty"`
	Replies   []Reply   `json:"replies,omitempty" datastore:",noindex"`
	Upvotes   int       `json:"upvotes,omitempty"`
	Downvotes int       `json:"downvotes,omitempty"`
	ModuleID  int64     `json:"module_id,omitempty"`
}

type ThreadRepository interface {
	CreateThread(Thread *Thread) (*datastore.Key, error)
	GetAllThreads() ([]*Thread, error)
	DeleteThread(id int64) error
	GetThreadByID(id int64) (*Thread, error)
	UpdateThread(id int64, Thread *Thread) (*datastore.Key, error)
	GetAllThreadsByModuleID(moduleID int64) ([]*Thread, error)
}
