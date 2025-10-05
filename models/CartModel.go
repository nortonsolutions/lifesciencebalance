package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// CartItem represents a single item in a shopping cart
type CartItem struct {
	ProductID int64 `json:"product_id,omitempty"`
	Quantity  int   `json:"quantity,omitempty"`
}

// Cart model - new for shopping cart functionality
type Cart struct {
	KeyID      int64      `json:"id"`
	CustomerID int64      `json:"customer_id,omitempty"`
	Items      []CartItem `json:"items,omitempty" datastore:",noindex"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
}

type CartRepository interface {
	CreateCart(cart *Cart) (*datastore.Key, error)
	GetCartByCustomerID(customerID int64) (*Cart, error)
	UpdateCart(id int64, cart *Cart) (*datastore.Key, error)
	DeleteCart(id int64) error
	ClearCart(customerID int64) error
}
