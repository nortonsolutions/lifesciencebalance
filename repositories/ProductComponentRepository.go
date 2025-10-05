package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type ProductComponentRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewProductComponentRepository(client *datastore.Client, ctx context.Context) models.ProductComponentRepository {
	return &ProductComponentRepositoryImpl{client: client, ctx: ctx}
}

func (r *ProductComponentRepositoryImpl) CreateProductComponent(component *models.ProductComponent) (*datastore.Key, error) {
	key := datastore.IncompleteKey("ProductComponent", nil)
	key, err := r.client.Put(r.ctx, key, component)
	if err != nil {
		return nil, err
	}
	component.KeyID = key.ID
	return key, nil
}

func (r *ProductComponentRepositoryImpl) GetAllProductComponents() ([]*models.ProductComponent, error) {
	var components []*models.ProductComponent
	query := datastore.NewQuery("ProductComponent")
	keys, err := r.client.GetAll(r.ctx, query, &components)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		components[i].KeyID = key.ID
	}
	return components, nil
}

func (r *ProductComponentRepositoryImpl) DeleteProductComponent(id int64) error {
	key := datastore.IDKey("ProductComponent", id, nil)
	return r.client.Delete(r.ctx, key)
}

func (r *ProductComponentRepositoryImpl) GetProductComponentByID(id int64) (*models.ProductComponent, error) {
	key := datastore.IDKey("ProductComponent", id, nil)
	var component models.ProductComponent
	if err := r.client.Get(r.ctx, key, &component); err != nil {
		return nil, err
	}
	component.KeyID = key.ID
	return &component, nil
}

func (r *ProductComponentRepositoryImpl) UpdateProductComponent(id int64, component *models.ProductComponent) (*datastore.Key, error) {
	key := datastore.IDKey("ProductComponent", id, nil)
	key, err := r.client.Put(r.ctx, key, component)
	if err != nil {
		return nil, err
	}
	component.KeyID = key.ID
	return key, nil
}

func (r *ProductComponentRepositoryImpl) GetAllProductComponentsByProductID(productID int64) ([]*models.ProductComponent, error) {
	var components []*models.ProductComponent
	query := datastore.NewQuery("ProductComponent").Filter("ProductID =", productID)
	keys, err := r.client.GetAll(r.ctx, query, &components)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		components[i].KeyID = key.ID
	}
	return components, nil
}
