package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newThreadRepository
func NewThreadRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Thread
func (r *BaseRepository) CreateThread(Thread *models.Thread) (*datastore.Key, error) {
	var pk1 *datastore.PendingKey
	var err error

	commit, err := r.client.RunInTransaction(r.ctx, func(tx *datastore.Transaction) error {
		pk1, err = tx.Put(datastore.IncompleteKey("Thread", nil), Thread)
		if err != nil {
			return err
		}
		return nil
	})

	k1 := commit.Key(pk1)
	return k1, err
}

// GetAllThreads returns all Threads
func (r *BaseRepository) GetAllThreads() ([]*models.Thread, error) {
	var Threads []*models.Thread
	query := datastore.NewQuery("Thread")
	keys, err := r.client.GetAll(r.ctx, query, &Threads)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Threads[i].KeyID = key.ID
	}

	return Threads, nil
}

// getThreadByID returns a Thread by id
func (r *BaseRepository) GetThreadByID(id int64) (*models.Thread, error) {
	Thread := new(models.Thread)

	k := datastore.IDKey("Thread", id, nil)
	if err := r.client.Get(r.ctx, k, Thread); err != nil {
		return nil, err
	}

	return Thread, nil
}

// UpdateThread updates a Thread
func (r *BaseRepository) UpdateThread(id int64, Thread *models.Thread) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Thread", id, nil), Thread)
}

// DeleteThread deletes a Thread
func (r *BaseRepository) DeleteThread(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Thread", id, nil))
}

// GetAllThreadsByModuleID
func (r *BaseRepository) GetAllThreadsByModuleID(moduleID int64) ([]*models.Thread, error) {
	var Threads []*models.Thread
	query := datastore.NewQuery("Thread").FilterField("ModuleID", "=", moduleID)
	keys, err := r.client.GetAll(r.ctx, query, &Threads)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Threads[i].KeyID = key.ID
	}

	return Threads, nil
}
