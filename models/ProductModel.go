package models

import "cloud.google.com/go/datastore"

// Product model - evolved from Course for e-commerce
type Product struct {
	KeyID       int64   `json:"id"` //gorm:"primary_key,autoIncrement"
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty" datastore:",noindex"`
	Price       float64 `json:"price,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	Category    string  `json:"category,omitempty"`
	VendorID    int64   `json:"vendor_id,omitempty"`
	Components  []int64 `json:"components,omitempty" datastore:",noindex"` // was Modules
	Approved    bool    `json:"approved"`
	InStock     bool    `json:"in_stock"`
	SKU         string  `json:"sku,omitempty"`
	Tags        []string `json:"tags,omitempty" datastore:",noindex"`
}

type ProductRepository interface {
	CreateProduct(product *Product) (*datastore.Key, error)
	GetAllProducts() ([]*Product, error)
	DeleteProduct(id int64) error
	GetProductByID(id int64) (*Product, error)
	UpdateProduct(id int64, product *Product) (*datastore.Key, error)
	GetApprovedProducts() ([]*Product, error)
	GetUnapprovedProducts() ([]*Product, error)
	GetProductsByCategory(category string) ([]*Product, error)
	GetProductsByVendor(vendorID int64) ([]*Product, error)
}
