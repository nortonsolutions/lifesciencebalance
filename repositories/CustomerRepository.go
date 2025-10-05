package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type CustomerRepositoryImpl struct {
	client *datastore.Client
	ctx    context.Context
}

func NewCustomerRepository(client *datastore.Client, ctx context.Context) models.CustomerRepository {
	return &CustomerRepositoryImpl{client: client, ctx: ctx}
}

func (r *CustomerRepositoryImpl) CreateCustomer(customer *models.Customer) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Customer", nil)
	key, err := r.client.Put(r.ctx, key, customer)
	if err != nil {
		return nil, err
	}
	customer.KeyID = key.ID
	return key, nil
}

func (r *CustomerRepositoryImpl) GetAllCustomers() ([]*models.Customer, error) {
	var customers []*models.Customer
	query := datastore.NewQuery("Customer")
	keys, err := r.client.GetAll(r.ctx, query, &customers)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		customers[i].KeyID = key.ID
	}
	return customers, nil
}

func (r *CustomerRepositoryImpl) GetCustomerByID(id int64) (*models.Customer, error) {
	key := datastore.IDKey("Customer", id, nil)
	var customer models.Customer
	if err := r.client.Get(r.ctx, key, &customer); err != nil {
		return nil, err
	}
	customer.KeyID = key.ID
	return &customer, nil
}

func (r *CustomerRepositoryImpl) GetCustomerByUserID(userID int64) (*models.Customer, error) {
	var customers []*models.Customer
	query := datastore.NewQuery("Customer").Filter("UserID =", userID).Limit(1)
	keys, err := r.client.GetAll(r.ctx, query, &customers)
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, datastore.ErrNoSuchEntity
	}
	customers[0].KeyID = keys[0].ID
	return customers[0], nil
}

func (r *CustomerRepositoryImpl) UpdateCustomer(id int64, customer *models.Customer) (*datastore.Key, error) {
	key := datastore.IDKey("Customer", id, nil)
	key, err := r.client.Put(r.ctx, key, customer)
	if err != nil {
		return nil, err
	}
	customer.KeyID = key.ID
	return key, nil
}

func (r *CustomerRepositoryImpl) DeleteCustomer(id int64) error {
	key := datastore.IDKey("Customer", id, nil)
	return r.client.Delete(r.ctx, key)
}
