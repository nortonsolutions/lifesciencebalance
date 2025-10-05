package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// OrderItem represents a single item in an order
type OrderItem struct {
	ProductID int64   `json:"product_id,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	Price     float64 `json:"price,omitempty"` // Price at time of purchase
}

// Order model - new for e-commerce transactions
type Order struct {
	KeyID           int64       `json:"id"`
	CustomerID      int64       `json:"customer_id,omitempty"`
	Items           []OrderItem `json:"items,omitempty" datastore:",noindex"`
	TotalAmount     float64     `json:"total_amount,omitempty"`
	Status          string      `json:"status,omitempty"` // pending, completed, shipped, cancelled
	ShippingAddress string      `json:"shipping_address,omitempty" datastore:",noindex"`
	BillingAddress  string      `json:"billing_address,omitempty" datastore:",noindex"`
	PaymentMethod   string      `json:"payment_method,omitempty"`
	PaymentStatus   string      `json:"payment_status,omitempty"` // pending, completed, failed
	TrackingNumber  string      `json:"tracking_number,omitempty"`
	CreatedAt       time.Time   `json:"created_at,omitempty"`
	UpdatedAt       time.Time   `json:"updated_at,omitempty"`
}

type OrderRepository interface {
	CreateOrder(order *Order) (*datastore.Key, error)
	GetAllOrders() ([]*Order, error)
	GetOrderByID(id int64) (*Order, error)
	GetOrdersByCustomerID(customerID int64) ([]*Order, error)
	UpdateOrder(id int64, order *Order) (*datastore.Key, error)
	DeleteOrder(id int64) error
	GetOrdersByStatus(status string) ([]*Order, error)
}
