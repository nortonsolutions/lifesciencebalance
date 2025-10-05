package models

import "cloud.google.com/go/datastore"

// create Role model
// Common roles: Student, Teacher, Admin, Customer, Vendor
type Role struct {
	// auto increment id
	KeyID        int64  `json:"id"` //gorm:"primary_key,autoIncrement"
	Name         string `json:"name,omitempty"`
	NumericValue int    `json:"numeric_value,omitempty"`
}

// RoleRepository ..
type RoleRepository interface {
	CreateRole(Role *Role) (*datastore.Key, error)
	GetAllRoles() ([]*Role, error)
	DeleteRole(id int64) error
	GetRoleByID(id int64) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	UpdateRole(id int64, Role *Role) (*datastore.Key, error)
}
