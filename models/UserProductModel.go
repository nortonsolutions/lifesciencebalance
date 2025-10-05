package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// JSON date string will look like this: "2021-02-18T21:54:42.123Z"
// err := time.Parse(time.RFC3339, dateString) //RFC 3339 is a profile for ISO 8601

type UserProduct struct {
	KeyID       int64     `json:"id"` //gorm:"primary_key,autoIncrement"
	UserID      int64     `json:"user_id,omitempty"`
	ProductID   int64     `json:"product_id,omitempty"`
	Grade       int       `json:"grade,omitempty"`
	StartedOn   time.Time `json:"started_on,omitempty"`
	CompletedOn time.Time `json:"completed_on,omitempty"`
	Role        string    `json:"role,omitempty"` // "customer" or "vendor"
}

type UserProductRepository interface {
	CreateUserProduct(userProduct *UserProduct) (*datastore.Key, error)
	GetAllUserProducts() ([]*UserProduct, error)
	DeleteUserProduct(id int64) error
	GetUserProductsByUserID(userID int64) ([]*UserProduct, error)
	GetUserProductsByProductID(productID int64) ([]*UserProduct, error)
	GetUserProductByUserIDAndProductID(userID int64, productID int64) (*UserProduct, error)
	UpdateUserProduct(id int64, userProduct *UserProduct) (*datastore.Key, error)
	GetProductsByUserID(userID int64) ([]*Product, error)
	GetUsersByProductID(productID int64) ([]*User, error)
	GetVendorsByProductID(productID int64) ([]*User, error)
}
