package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newProjectRepository
func NewProjectRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Project
func (r *BaseRepository) CreateProject(Project *models.Project) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("Project", nil), Project)
}

// GetAllProjects returns all Projects
func (r *BaseRepository) GetAllProjects() ([]*models.Project, error) {
	var Projects []*models.Project
	query := datastore.NewQuery("Project")
	keys, err := r.client.GetAll(r.ctx, query, &Projects)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Projects[i].KeyID = key.ID
	}

	return Projects, nil
}

// getProjectByID returns a Project by id
func (r *BaseRepository) GetProjectByID(id int64) (*models.Project, error) {
	Project := new(models.Project)

	k := datastore.IDKey("Project", id, nil)
	if err := r.client.Get(r.ctx, k, Project); err != nil {
		return nil, err
	}

	return Project, nil
}

// UpdateProject updates a Project
func (r *BaseRepository) UpdateProject(id int64, Project *models.Project) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Project", id, nil), Project)
}

// DeleteProject deletes a Project
func (r *BaseRepository) DeleteProject(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Project", id, nil))
}

func (r *BaseRepository) GetProjectsByUserID(userID int64) ([]*models.Project, error) {
	var Projects []*models.Project
	query := datastore.NewQuery("Project").FilterField("UserID", "=", userID)
	_, err := r.client.GetAll(r.ctx, query, &Projects)
	if err != nil {
		return nil, err
	}

	return Projects, nil
}

func (r *BaseRepository) GetProjectsByCourseID(courseID int64) ([]*models.Project, error) {
	var Projects []*models.Project
	query := datastore.NewQuery("Project").FilterField("CourseID", "=", courseID)
	_, err := r.client.GetAll(r.ctx, query, &Projects)
	if err != nil {
		return nil, err
	}

	return Projects, nil
}

func (r *BaseRepository) GetProjectsByModuleID(moduleID int64) ([]*models.Project, error) {
	var Projects []*models.Project
	query := datastore.NewQuery("Project").FilterField("ModuleID", "=", moduleID)
	_, err := r.client.GetAll(r.ctx, query, &Projects)
	if err != nil {
		return nil, err
	}

	return Projects, nil
}

func (r *BaseRepository) GetProjectByUserIDandModuleID(userID int64, moduleID int64) (*models.Project, error) {
	Project := new(models.Project)

	query := datastore.NewQuery("Project").FilterField("UserID", "=", userID).FilterField("ModuleID", "=", moduleID)
	_, err := r.client.GetAll(r.ctx, query, &Project)
	if err != nil {
		return nil, err
	}

	return Project, nil
}
