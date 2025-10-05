package models

import "cloud.google.com/go/datastore"

type Route struct {
	KeyID           int64  `json:"id"` //gorm:"primary_key,autoIncrement"
	Name            string `json:"name,omitempty"`
	PermissionLevel int    `json:"numeric_value,omitempty"`
}

type RouteRepository interface {
	CreateRoute(Route *Route) (*datastore.Key, error)
	GetAllRoutes() ([]*Route, error)
	DeleteRoute(id int64) error
	GetRouteByID(id int64) (*Route, error)
	GetRouteByName(name string) (*Route, error)
	UpdateRoute(id int64, Route *Route) (*datastore.Key, error)
}
