package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newRouteRepository
func NewRouteRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Route
func (r *BaseRepository) CreateRoute(Route *models.Route) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("Route", nil), Route)
}

// GetAllRoutes returns all Routes
func (r *BaseRepository) GetAllRoutes() ([]*models.Route, error) {
	var Routes []*models.Route
	query := datastore.NewQuery("Route")
	keys, err := r.client.GetAll(r.ctx, query, &Routes)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Routes[i].KeyID = key.ID
	}

	return Routes, nil
}

// getRouteByID returns a Route by id
func (r *BaseRepository) GetRouteByID(id int64) (*models.Route, error) {
	Route := new(models.Route)

	k := datastore.IDKey("Route", id, nil)
	if err := r.client.Get(r.ctx, k, Route); err != nil {
		return nil, err
	}

	return Route, nil
}

// getRouteByRoutename returns a Route by Routename
func (r *BaseRepository) GetRouteByName(Routename string) (*models.Route, error) {
	Route := new(models.Route)
	query := datastore.NewQuery("Route").
		FilterField("Name", "=", Routename)

	key, err := r.client.Run(r.ctx, query).Next(Route)
	if err != nil {
		return nil, err
	}

	if err := r.client.Get(r.ctx, key, Route); err != nil {
		return nil, err
	}

	return Route, nil
}

// UpdateRoute updates a Route
func (r *BaseRepository) UpdateRoute(id int64, Route *models.Route) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Route", id, nil), Route)
}

// DeleteRoute deletes a Route
func (r *BaseRepository) DeleteRoute(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Route", id, nil))
}
