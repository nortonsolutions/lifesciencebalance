package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type ProductRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewProductRepository(client *datastore.Client, ctx context.Context) models.ProductRepository {
	return &ProductRepositoryImpl{client: client, ctx: ctx}
}

func (r *ProductRepositoryImpl) CreateProduct(product *models.Product) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Product", nil)
	key, err := r.client.Put(r.ctx, key, product)
	if err != nil {
		return nil, err
	}
	product.KeyID = key.ID
	return key, nil
}

func (r *ProductRepositoryImpl) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	query := datastore.NewQuery("Product")
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		products[i].KeyID = key.ID
	}
	return products, nil
}

func (r *ProductRepositoryImpl) DeleteProduct(id int64) error {
	key := datastore.IDKey("Product", id, nil)
	return r.client.Delete(r.ctx, key)
}

func (r *ProductRepositoryImpl) GetProductByID(id int64) (*models.Product, error) {
	key := datastore.IDKey("Product", id, nil)
	var product models.Product
	if err := r.client.Get(r.ctx, key, &product); err != nil {
		return nil, err
	}
	product.KeyID = key.ID
	return &product, nil
}

func (r *ProductRepositoryImpl) UpdateProduct(id int64, product *models.Product) (*datastore.Key, error) {
	key := datastore.IDKey("Product", id, nil)
	key, err := r.client.Put(r.ctx, key, product)
	if err != nil {
		return nil, err
	}
	product.KeyID = key.ID
	return key, nil
}

func (r *ProductRepositoryImpl) GetApprovedProducts() ([]*models.Product, error) {
	var products []*models.Product
	query := datastore.NewQuery("Product").Filter("Approved =", true)
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		products[i].KeyID = key.ID
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetUnapprovedProducts() ([]*models.Product, error) {
	var products []*models.Product
	query := datastore.NewQuery("Product").Filter("Approved =", false)
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		products[i].KeyID = key.ID
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetProductsByCategory(category string) ([]*models.Product, error) {
	var products []*models.Product
	query := datastore.NewQuery("Product").Filter("Category =", category)
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		products[i].KeyID = key.ID
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetProductsByVendor(vendorID int64) ([]*models.Product, error) {
	var products []*models.Product
	query := datastore.NewQuery("Product").Filter("VendorID =", vendorID)
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		products[i].KeyID = key.ID
	}
	return products, nil
}
