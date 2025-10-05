package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type CartRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewCartRepository(client *datastore.Client, ctx context.Context) models.CartRepository {
	return &CartRepositoryImpl{client: client, ctx: ctx}
}

func (r *CartRepositoryImpl) CreateCart(cart *models.Cart) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Cart", nil)
	key, err := r.client.Put(r.ctx, key, cart)
	if err != nil {
		return nil, err
	}
	cart.KeyID = key.ID
	return key, nil
}

func (r *CartRepositoryImpl) GetCartByCustomerID(customerID int64) (*models.Cart, error) {
	var carts []*models.Cart
	query := datastore.NewQuery("Cart").Filter("CustomerID =", customerID).Limit(1)
	keys, err := r.client.GetAll(r.ctx, query, &carts)
	if err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}
	carts[0].KeyID = keys[0].ID
	return carts[0], nil
}

func (r *CartRepositoryImpl) UpdateCart(id int64, cart *models.Cart) (*datastore.Key, error) {
	key := datastore.IDKey("Cart", id, nil)
	key, err := r.client.Put(r.ctx, key, cart)
	if err != nil {
		return nil, err
	}
	cart.KeyID = key.ID
	return key, nil
}

func (r *CartRepositoryImpl) DeleteCart(id int64) error {
	key := datastore.IDKey("Cart", id, nil)
	return r.client.Delete(r.ctx, key)
}

func (r *CartRepositoryImpl) ClearCart(customerID int64) error {
	cart, err := r.GetCartByCustomerID(customerID)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil // Cart doesn't exist, nothing to clear
		}
		return err
	}
	cart.Items = []models.CartItem{}
	_, err = r.UpdateCart(cart.KeyID, cart)
	return err
}
