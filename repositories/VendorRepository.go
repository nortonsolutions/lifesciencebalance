package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type VendorRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewVendorRepository(client *datastore.Client, ctx context.Context) models.VendorRepository {
	return &VendorRepositoryImpl{client: client, ctx: ctx}
}

func (r *VendorRepositoryImpl) CreateVendor(vendor *models.Vendor) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Vendor", nil)
	key, err := r.client.Put(r.ctx, key, vendor)
	if err != nil {
		return nil, err
	}
	vendor.KeyID = key.ID
	return key, nil
}

func (r *VendorRepositoryImpl) GetAllVendors() ([]*models.Vendor, error) {
	var vendors []*models.Vendor
	query := datastore.NewQuery("Vendor")
	keys, err := r.client.GetAll(r.ctx, query, &vendors)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		vendors[i].KeyID = key.ID
	}
	return vendors, nil
}

func (r *VendorRepositoryImpl) GetVendorByID(id int64) (*models.Vendor, error) {
	key := datastore.IDKey("Vendor", id, nil)
	var vendor models.Vendor
	if err := r.client.Get(r.ctx, key, &vendor); err != nil {
		return nil, err
	}
	vendor.KeyID = key.ID
	return &vendor, nil
}

func (r *VendorRepositoryImpl) GetVendorByUserID(userID int64) (*models.Vendor, error) {
	var vendors []*models.Vendor
	query := datastore.NewQuery("Vendor").Filter("UserID =", userID).Limit(1)
	keys, err := r.client.GetAll(r.ctx, query, &vendors)
	if err != nil {
		return nil, err
	}
	if len(vendors) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}
	vendors[0].KeyID = keys[0].ID
	return vendors[0], nil
}

func (r *VendorRepositoryImpl) UpdateVendor(id int64, vendor *models.Vendor) (*datastore.Key, error) {
	key := datastore.IDKey("Vendor", id, nil)
	key, err := r.client.Put(r.ctx, key, vendor)
	if err != nil {
		return nil, err
	}
	vendor.KeyID = key.ID
	return key, nil
}

func (r *VendorRepositoryImpl) DeleteVendor(id int64) error {
	key := datastore.IDKey("Vendor", id, nil)
	return r.client.Delete(r.ctx, key)
}

func (r *VendorRepositoryImpl) GetApprovedVendors() ([]*models.Vendor, error) {
	var vendors []*models.Vendor
	query := datastore.NewQuery("Vendor").Filter("Approved =", true)
	keys, err := r.client.GetAll(r.ctx, query, &vendors)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		vendors[i].KeyID = key.ID
	}
	return vendors, nil
}
