package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newCourseRepository
func NewCourseRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new Course
func (r *BaseRepository) CreateCourse(Course *models.Course) (*datastore.Key, error) {
	var pk1 *datastore.PendingKey
	var err error

	commit, err := r.client.RunInTransaction(r.ctx, func(tx *datastore.Transaction) error {
		pk1, err = tx.Put(datastore.IncompleteKey("Course", nil), Course)
		if err != nil {
			return err
		}
		return nil
	})

	k1 := commit.Key(pk1)
	return k1, err
}

// GetAllCourses returns all Courses
func (r *BaseRepository) GetAllCourses() ([]*models.Course, error) {
	var Courses []*models.Course
	query := datastore.NewQuery("Course")
	keys, err := r.client.GetAll(r.ctx, query, &Courses)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		Courses[i].KeyID = key.ID
	}

	return Courses, nil
}

// getCourseByID returns a Course by id
func (r *BaseRepository) GetCourseByID(id int64) (*models.Course, error) {
	Course := new(models.Course)

	k := datastore.IDKey("Course", id, nil)
	if err := r.client.Get(r.ctx, k, Course); err != nil {
		return nil, err
	}

	Course.KeyID = k.ID
	return Course, nil
}

// UpdateCourse updates a Course
func (r *BaseRepository) UpdateCourse(id int64, Course *models.Course) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("Course", id, nil), Course)
}

// DeleteCourse deletes a Course
func (r *BaseRepository) DeleteCourse(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("Course", id, nil))
}

// GetApprovedCourses returns only approved courses
func (r *BaseRepository) GetApprovedCourses() ([]*models.Course, error) {
	var courses []*models.Course
	// creates a new query for approved courses
	query := datastore.NewQuery("Course").Filter("Approved =", true)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &courses)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each course
	for i, key := range keys {
		courses[i].KeyID = key.ID
	}

	return courses, nil
}

// GetUnapprovedCourses returns only unapproved courses
func (r *BaseRepository) GetUnapprovedCourses() ([]*models.Course, error) {
	var courses []*models.Course
	// creates a new query for unapproved courses
	query := datastore.NewQuery("Course").Filter("Approved =", false)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &courses)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each course
	for i, key := range keys {
		courses[i].KeyID = key.ID
	}

	return courses, nil
}

// GetCoursesByDepartment returns courses filtered by department
func (r *BaseRepository) GetCoursesByDepartment(department string) ([]*models.Course, error) {
	var courses []*models.Course
	// creates a new query for courses in a specific department
	query := datastore.NewQuery("Course").Filter("Department =", department)

	// Run the query
	keys, err := r.client.GetAll(r.ctx, query, &courses)
	if err != nil {
		return nil, err
	}

	// Sets the key ID for each course
	for i, key := range keys {
		courses[i].KeyID = key.ID
	}

	return courses, nil
}
