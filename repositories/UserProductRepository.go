package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newUserProductRepository
func NewUserProductRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new UserProduct
func (r *BaseRepository) CreateUserProduct(UserProduct *models.UserProduct) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("UserProduct", nil), UserProduct)
}

// GetAllUserProducts returns all UserProducts
func (r *BaseRepository) GetAllUserProducts() ([]*models.UserProduct, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct")
	keys, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		UserProducts[i].KeyID = key.ID
	}

	return UserProducts, nil
}

// DeleteUserProduct deletes a UserProduct
func (r *BaseRepository) DeleteUserProduct(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("UserProduct", id, nil))
}

func (r *BaseRepository) UpdateUserProduct(id int64, UserProduct *models.UserProduct) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("UserProduct", id, nil), UserProduct)
}

func (r *BaseRepository) GetUserProductsByUserID(userID int64) ([]*models.UserProduct, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("UserID", "=", userID)
	_, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	return UserProducts, nil
}

func (r *BaseRepository) GetUserProductsByProductID(productID int64) ([]*models.UserProduct, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("ProductID", "=", productID)
	_, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	return UserProducts, nil
}

func (r *BaseRepository) GetUserProductByUserIDAndProductID(userID int64, productID int64) (*models.UserProduct, error) {

	// TODO: currently using an array but this should only return one
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("UserID", "=", userID).FilterField("ProductID", "=", productID)
	keys, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	// add key values to UserProducts
	for i, key := range keys {
		UserProducts[i].KeyID = key.ID
	}

	return UserProducts[0], nil
}

func (r *BaseRepository) GetProductsByUserID(userID int64) ([]*models.Product, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("UserID", "=", userID)
	_, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	var Products []*models.Product
	for _, UserProduct := range UserProducts {
		Product, err := r.GetProductByID(UserProduct.ProductID)
		if err != nil {
			return nil, err
		}

		Product.KeyID = UserProduct.ProductID
		Products = append(Products, Product)
	}

	return Products, nil
}

func (r *BaseRepository) GetUsersByProductID(productID int64) ([]*models.User, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("ProductID", "=", productID).FilterField("Role", "!=", "vendor")
	_, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	var Users []*models.User
	for _, UserProduct := range UserProducts {
		User, err := r.GetUserByID(UserProduct.UserID)
		if err != nil {
			return nil, err
		}
		Users = append(Users, User)
	}

	return Users, nil
}

func (r *BaseRepository) GetVendorsByProductID(productID int64) ([]*models.User, error) {
	var UserProducts []*models.UserProduct
	query := datastore.NewQuery("UserProduct").FilterField("ProductID", "=", productID).FilterField("Role", "=", "vendor")
	_, err := r.client.GetAll(r.ctx, query, &UserProducts)
	if err != nil {
		return nil, err
	}

	var Users []*models.User
	for _, UserProduct := range UserProducts {
		User, err := r.GetUserByID(UserProduct.UserID)
		if err != nil {
			return nil, err
		}
		Users = append(Users, User)
	}

	return Users, nil
}
