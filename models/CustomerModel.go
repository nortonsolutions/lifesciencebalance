package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// Customer model - evolved from User with Student role context
type Customer struct {
	KeyID            int64     `json:"id"`
	UserID           int64     `json:"user_id,omitempty"` // Reference to User
	ShippingAddress  string    `json:"shipping_address,omitempty" datastore:",noindex"`
	BillingAddress   string    `json:"billing_address,omitempty" datastore:",noindex"`
	Phone            string    `json:"phone,omitempty"`
	PreferredPayment string    `json:"preferred_payment,omitempty"`
	CreatedOn        time.Time `json:"created_on,omitempty"`
}

type CustomerRepository interface {
	CreateCustomer(customer *Customer) (*datastore.Key, error)
	GetAllCustomers() ([]*Customer, error)
	GetCustomerByID(id int64) (*Customer, error)
	GetCustomerByUserID(userID int64) (*Customer, error)
	UpdateCustomer(id int64, customer *Customer) (*datastore.Key, error)
	DeleteCustomer(id int64) error
}
