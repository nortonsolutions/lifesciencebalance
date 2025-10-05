package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newModuleRepository
func NewModuleRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Module and return the committed key
func (r *BaseRepository) CreateModule(Module *models.Module) (*datastore.Key, error) {
	var pk1 *datastore.PendingKey
	var err error

	commit, err := r.client.RunInTransaction(r.ctx, func(tx *datastore.Transaction) error {
		pk1, err = tx.Put(datastore.IncompleteKey("Module", nil), Module)
		if err != nil {
			return err
		}
		return nil
	})

	k1 := commit.Key(pk1)
	return k1, err
}

// GetAllModules returns all Modules
func (r *BaseRepository) GetAllModules() ([]*models.Module, error) {
	var Modules []*models.Module
	query := datastore.NewQuery("Module")
	keys, err := r.client.GetAll(r.ctx, query, &Modules)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Modules[i].KeyID = key.ID
	}

	return Modules, nil
}

// getModuleByID returns a Module by id
func (r *BaseRepository) GetModuleByID(id int64) (*models.Module, error) {
	Module := new(models.Module)

	k := datastore.IDKey("Module", id, nil)
	if err := r.client.Get(r.ctx, k, Module); err != nil {
		return nil, err
	}

	Module.KeyID = k.ID

	return Module, nil
}

// UpdateModule updates a Module
func (r *BaseRepository) UpdateModule(id int64, Module *models.Module) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Module", id, nil), Module)
}

// DeleteModule deletes a Module
func (r *BaseRepository) DeleteModule(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Module", id, nil))
}

func (r *BaseRepository) GetAllModulesByProductID(productID int64) ([]*models.Module, error) {
	var Modules []*models.Module
	query := datastore.NewQuery("Module").FilterField("ProductID", "=", productID)
	keys, err := r.client.GetAll(r.ctx, query, &Modules)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Modules[i].KeyID = key.ID
	}

	return Modules, nil
}
