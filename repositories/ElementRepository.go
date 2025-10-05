package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newElementRepository
func NewElementRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Element
func (r *BaseRepository) CreateElement(Element *models.Element) (*datastore.Key, error) {
	var pk1 *datastore.PendingKey
	var err error

	commit, err := r.client.RunInTransaction(r.ctx, func(tx *datastore.Transaction) error {
		pk1, err = tx.Put(datastore.IncompleteKey("Element", nil), Element)
		if err != nil {
			return err
		}
		return nil
	})

	k1 := commit.Key(pk1)
	return k1, err
}

// GetAllElements returns all Elements
func (r *BaseRepository) GetAllElements() ([]*models.Element, error) {
	var Elements []*models.Element
	query := datastore.NewQuery("Element")
	keys, err := r.client.GetAll(r.ctx, query, &Elements)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Elements[i].KeyID = key.ID
	}

	return Elements, nil
}

// getElementByID returns a Element by id
func (r *BaseRepository) GetElementByID(id int64) (*models.Element, error) {
	Element := new(models.Element)

	k := datastore.IDKey("Element", id, nil)
	if err := r.client.Get(r.ctx, k, Element); err != nil {
		return nil, err
	}

	Element.KeyID = k.ID
	return Element, nil
}

// UpdateElement updates a Element
func (r *BaseRepository) UpdateElement(id int64, Element *models.Element) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Element", id, nil), Element)
}

// DeleteElement deletes a Element
func (r *BaseRepository) DeleteElement(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Element", id, nil))
}
