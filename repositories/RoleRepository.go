package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newRoleRepository
func NewRoleRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Role
func (r *BaseRepository) CreateRole(Role *models.Role) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("Role", nil), Role)
}

// GetAllRoles returns all Roles
func (r *BaseRepository) GetAllRoles() ([]*models.Role, error) {
	var Roles []*models.Role
	query := datastore.NewQuery("Role")
	keys, err := r.client.GetAll(r.ctx, query, &Roles)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Roles[i].KeyID = key.ID
	}

	return Roles, nil
}

// getRoleByID returns a Role by id
func (r *BaseRepository) GetRoleByID(id int64) (*models.Role, error) {
	Role := new(models.Role)

	k := datastore.IDKey("Role", id, nil)
	if err := r.client.Get(r.ctx, k, Role); err != nil {
		return nil, err
	}

	return Role, nil
}

// getRoleByRolename returns a Role by Rolename
func (r *BaseRepository) GetRoleByName(Rolename string) (*models.Role, error) {
	Role := new(models.Role)
	query := datastore.NewQuery("Role").
		FilterField("Name", "=", Rolename)

	key, err := r.client.Run(r.ctx, query).Next(Role)
	if err != nil {
		return nil, err
	}

	if err := r.client.Get(r.ctx, key, Role); err != nil {
		return nil, err
	}

	return Role, nil
}

// UpdateRole updates a Role
func (r *BaseRepository) UpdateRole(id int64, Role *models.Role) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Role", id, nil), Role)
}

// DeleteRole deletes a Role
func (r *BaseRepository) DeleteRole(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Role", id, nil))
}
