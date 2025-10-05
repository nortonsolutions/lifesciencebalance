package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newProductRepository
func NewProductRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Product
func (r *BaseRepository) CreateProduct(Product *models.Product) (*datastore.Key, error) {
	var pk1 *datastore.PendingKey
	var err error

	commit, err := r.client.RunInTransaction(r.ctx, func(tx *datastore.Transaction) error {
		pk1, err = tx.Put(datastore.IncompleteKey("Product", nil), Product)
		if err != nil {
			return err
		}
		return nil
	})

	k1 := commit.Key(pk1)
	return k1, err
}

// GetAllProducts returns all Products
func (r *BaseRepository) GetAllProducts() ([]*models.Product, error) {
	var Products []*models.Product
	query := datastore.NewQuery("Product")
	keys, err := r.client.GetAll(r.ctx, query, &Products)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Products[i].KeyID = key.ID
	}

	return Products, nil
}

// getProductByID returns a Product by id
func (r *BaseRepository) GetProductByID(id int64) (*models.Product, error) {
	Product := new(models.Product)

	k := datastore.IDKey("Product", id, nil)
	if err := r.client.Get(r.ctx, k, Product); err != nil {
		return nil, err
	}

	Product.KeyID = k.ID
	return Product, nil
}

// UpdateProduct updates a Product
func (r *BaseRepository) UpdateProduct(id int64, Product *models.Product) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Product", id, nil), Product)
}

// DeleteProduct deletes a Product
func (r *BaseRepository) DeleteProduct(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Product", id, nil))
}

// GetApprovedProducts returns only approved products
func (r *BaseRepository) GetApprovedProducts() ([]*models.Product, error) {
	var products []*models.Product
	// creates a new query for approved products
	query := datastore.NewQuery("Product").Filter("Approved =", true)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each product
	for i, key := range keys {
		products[i].KeyID = key.ID
	}

	return products, nil
}

// GetUnapprovedProducts returns only unapproved products
func (r *BaseRepository) GetUnapprovedProducts() ([]*models.Product, error) {
	var products []*models.Product
	// creates a new query for unapproved products
	query := datastore.NewQuery("Product").Filter("Approved =", false)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each product
	for i, key := range keys {
		products[i].KeyID = key.ID
	}

	return products, nil
}

// GetProductsByDepartment returns products filtered by department
func (r *BaseRepository) GetProductsByDepartment(department string) ([]*models.Product, error) {
	var products []*models.Product
	// creates a new query for products in a specific department
	query := datastore.NewQuery("Product").Filter("Department =", department)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &products)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each product
	for i, key := range keys {
		products[i].KeyID = key.ID
	}

	return products, nil
}
