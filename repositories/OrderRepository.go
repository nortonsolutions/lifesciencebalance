package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type OrderRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewOrderRepository(client *datastore.Client, ctx context.Context) models.OrderRepository {
	return &OrderRepositoryImpl{client: client, ctx: ctx}
}

func (r *OrderRepositoryImpl) CreateOrder(order *models.Order) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Order", nil)
	key, err := r.client.Put(r.ctx, key, order)
	if err != nil {
		return nil, err
	}
	order.KeyID = key.ID
	return key, nil
}

func (r *OrderRepositoryImpl) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order
	query := datastore.NewQuery("Order")
	keys, err := r.client.GetAll(r.ctx, query, &orders)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		orders[i].KeyID = key.ID
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) GetOrderByID(id int64) (*models.Order, error) {
	key := datastore.IDKey("Order", id, nil)
	var order models.Order
	if err := r.client.Get(r.ctx, key, &order); err != nil {
		return nil, err
	}
	order.KeyID = key.ID
	return &order, nil
}

func (r *OrderRepositoryImpl) GetOrdersByCustomerID(customerID int64) ([]*models.Order, error) {
	var orders []*models.Order
	query := datastore.NewQuery("Order").Filter("CustomerID =", customerID)
	keys, err := r.client.GetAll(r.ctx, query, &orders)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		orders[i].KeyID = key.ID
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) UpdateOrder(id int64, order *models.Order) (*datastore.Key, error) {
	key := datastore.IDKey("Order", id, nil)
	key, err := r.client.Put(r.ctx, key, order)
	if err != nil {
		return nil, err
	}
	order.KeyID = key.ID
	return key, nil
}

func (r *OrderRepositoryImpl) DeleteOrder(id int64) error {
	key := datastore.IDKey("Order", id, nil)
	return r.client.Delete(r.ctx, key)
}

func (r *OrderRepositoryImpl) GetOrdersByStatus(status string) ([]*models.Order, error) {
	var orders []*models.Order
	query := datastore.NewQuery("Order").Filter("Status =", status)
	keys, err := r.client.GetAll(r.ctx, query, &orders)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		orders[i].KeyID = key.ID
	}
	return orders, nil
}
