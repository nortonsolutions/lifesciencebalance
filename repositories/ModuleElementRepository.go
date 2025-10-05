package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newModuleElementRepository
func NewModuleElementRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new ModuleElement
func (r *BaseRepository) CreateModuleElement(ModuleElement *models.ModuleElement) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("ModuleElement", nil), ModuleElement)
}

// GetAllModuleElements returns all ModuleElements
func (r *BaseRepository) GetAllModuleElements() ([]*models.ModuleElement, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement")
	keys, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		ModuleElements[i].KeyID = key.ID
	}

	return ModuleElements, nil
}

// DeleteModuleElement deletes a ModuleElement
func (r *BaseRepository) DeleteModuleElement(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("ModuleElement", id, nil))
}

func (r *BaseRepository) UpdateModuleElement(id int64, ModuleElement *models.ModuleElement) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("ModuleElement", id, nil), ModuleElement)
}

func (r *BaseRepository) GetModuleElementsByModuleID(moduleID int64) ([]*models.ModuleElement, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("ModuleID", "=", moduleID)
	keys, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	// add key values to ModuleElements
	for i, key := range keys {
		ModuleElements[i].KeyID = key.ID
	}

	return ModuleElements, nil
}

func (r *BaseRepository) GetModuleElementsByElementID(elementID int64) ([]*models.ModuleElement, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("ElementID", "=", elementID)
	keys, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		ModuleElements[i].KeyID = key.ID
	}

	return ModuleElements, nil
}

func (r *BaseRepository) GetModuleElementByModuleIDAndElementID(moduleID int64, elementID int64) (*models.ModuleElement, error) {

	// TODO: currently using an array but this should only return one
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("ModuleID", "=", moduleID).FilterField("ElementID", "=", elementID)
	keys, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	// add key values to ModuleElements
	for i, key := range keys {
		ModuleElements[i].KeyID = key.ID
	}

	return ModuleElements[0], nil
}

func (r *BaseRepository) GetElementsByModuleID(moduleID int64) ([]*models.Element, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("ModuleID", "=", moduleID)
	_, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	var Elements []*models.Element
	for _, ModuleElement := range ModuleElements {
		Element, err := r.GetElementByID(ModuleElement.ElementID)
		if err != nil {
			return nil, err
		}

		Element.KeyID = ModuleElement.ElementID
		Elements = append(Elements, Element)
	}

	return Elements, nil
}

func (r *BaseRepository) GetModulesByElementID(elementID int64) ([]*models.Module, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("ElementID", "=", elementID)
	_, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	var Modules []*models.Module
	for _, ModuleElement := range ModuleElements {
		Module, err := r.GetModuleByID(ModuleElement.ModuleID)
		if err != nil {
			return nil, err
		}
		Modules = append(Modules, Module)
	}

	return Modules, nil
}

func (r *BaseRepository) GetElementsByInstructorID(elementID int64) ([]*models.Element, error) {
	var ModuleElements []*models.ModuleElement
	query := datastore.NewQuery("ModuleElement").FilterField("InstructorID", "=", elementID).FilterField("Role", "=", "instructor")
	_, err := r.client.GetAll(r.ctx, query, &ModuleElements)
	if err != nil {
		return nil, err
	}

	var Elements []*models.Element
	for _, ModuleElement := range ModuleElements {
		Element, err := r.GetElementByID(ModuleElement.ElementID)
		if err != nil {
			return nil, err
		}
		Elements = append(Elements, Element)
	}

	return Elements, nil
}
