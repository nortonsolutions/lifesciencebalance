package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

// Vendor model - evolved from User with Teacher role context
type Vendor struct {
	KeyID         int64     `json:"id"`
	UserID        int64     `json:"user_id,omitempty"` // Reference to User
	BusinessName  string    `json:"business_name,omitempty"`
	BusinessEmail string    `json:"business_email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Address       string    `json:"address,omitempty" datastore:",noindex"`
	TaxID         string    `json:"tax_id,omitempty"`
	PaymentInfo   string    `json:"payment_info,omitempty" datastore:",noindex"` // For receiving payments
	Approved      bool      `json:"approved"`
	CreatedOn     time.Time `json:"created_on,omitempty"`
}

type VendorRepository interface {
	CreateVendor(vendor *Vendor) (*datastore.Key, error)
	GetAllVendors() ([]*Vendor, error)
	GetVendorByID(id int64) (*Vendor, error)
	GetVendorByUserID(userID int64) (*Vendor, error)
	UpdateVendor(id int64, vendor *Vendor) (*datastore.Key, error)
	DeleteVendor(id int64) error
	GetApprovedVendors() ([]*Vendor, error)
}
