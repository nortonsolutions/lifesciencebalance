package models

import "cloud.google.com/go/datastore"

// ProductComponent model - evolved from Module for product page components
type ProductComponent struct {
	KeyID       int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty" datastore:",noindex"`
	Type        string `json:"type,omitempty"` // e.g., "image_gallery", "reviews", "specifications", "transaction"
	Content     string `json:"content,omitempty" datastore:",noindex"`
	SortKey     int    `json:"sort_key,omitempty"`
	ProductID   int64  `json:"product_id,omitempty"`
	VendorID    int64  `json:"vendor_id,omitempty"`
}

type ProductComponentRepository interface {
	CreateProductComponent(component *ProductComponent) (*datastore.Key, error)
	GetAllProductComponents() ([]*ProductComponent, error)
	DeleteProductComponent(id int64) error
	GetProductComponentByID(id int64) (*ProductComponent, error)
	UpdateProductComponent(id int64, component *ProductComponent) (*datastore.Key, error)
	GetAllProductComponentsByProductID(productID int64) ([]*ProductComponent, error)
}
