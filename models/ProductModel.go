package models

import "cloud.google.com/go/datastore"

type Product struct {
	KeyID       int64   `json:"id"` //gorm:"primary_key,autoIncrement"
	Name        string  `json:"name,omitempty"`
	HomeContent string  `json:"home_content,omitempty" datastore:",noindex"`
	Description string  `json:"description,omitempty" datastore:",noindex"`
	Modules     []int64 `json:"modules,omitempty" datastore:",noindex"`
	OwnerID     int64   `json:"owner_id,omitempty"`
	Approved    bool    `json:"approved"`
	Department  string  `json:"department,omitempty"`
}

type ProductRepository interface {
	CreateProduct(Product *Product) (*datastore.Key, error)
	GetAllProducts() ([]*Product, error)
	DeleteProduct(id int64) error
	GetProductByID(id int64) (*Product, error)
	UpdateProduct(id int64, Product *Product) (*datastore.Key, error)
	GetApprovedProducts() ([]*Product, error)
	GetUnapprovedProducts() ([]*Product, error)
	GetProductsByDepartment(department string) ([]*Product, error)
}
