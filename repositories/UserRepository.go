package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

type BaseRepository struct {
	client *datastore.Client
	ctx    context.Context
}

// newUserRepository
func NewUserRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new user
func (r *BaseRepository) CreateUser(user *models.User) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("User", nil), user)
}

// GetAllUsers returns all users
func (r *BaseRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	query := datastore.NewQuery("User")
	keys, err := r.client.GetAll(r.ctx, query, &users)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		users[i].KeyID = key.ID
	}

	return users, nil
}

// getUserByID returns a user by id
func (r *BaseRepository) GetUserByID(id int64) (*models.User, error) {
	user := new(models.User)

	k := datastore.IDKey("User", id, nil)
	if err := r.client.Get(r.ctx, k, user); err != nil {
		return nil, err
	}

	return user, nil
}

// getUserByUsername returns a user by username
func (r *BaseRepository) GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)
	query := datastore.NewQuery("User").
		FilterField("Username", "=", username)

	key, err := r.client.Run(r.ctx, query).Next(user)
	if err != nil {
		return nil, err
	}

	if err := r.client.Get(r.ctx, key, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user
func (r *BaseRepository) UpdateUser(id int64, user *models.User) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("User", id, nil), user)
}

// DeleteUser deletes a user
func (r *BaseRepository) DeleteUser(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("User", id, nil))
}

func (r *BaseRepository) GetUserByUsernameAndPassword(username string, password string) (*models.User, error) {

	user := new(models.User)
	query := datastore.NewQuery("User").
		FilterField("Username", "=", username).
		FilterField("Password", "=", password)

	// use client to run the query
	key, err := r.client.Run(r.ctx, query).Next(user)
	if err != nil {
		return nil, err
	}

	if err := r.client.Get(r.ctx, key, user); err != nil {
		return nil, err
	}

	return user, nil
}
