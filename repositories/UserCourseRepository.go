package repositories

import (
	"context"
	"restAPI/models"

	"cloud.google.com/go/datastore"
)

// newUserCourseRepository
func NewUserCourseRepository(client *datastore.Client, ctx context.Context) *BaseRepository {

	return &BaseRepository{
		client: client,
		ctx:    ctx,
	}
}

// Create a new UserCourse
func (r *BaseRepository) CreateUserCourse(UserCourse *models.UserCourse) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IncompleteKey("UserCourse", nil), UserCourse)
}

// GetAllUserCourses returns all UserCourses
func (r *BaseRepository) GetAllUserCourses() ([]*models.UserCourse, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse")
	keys, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		UserCourses[i].KeyID = key.ID
	}

	return UserCourses, nil
}

// DeleteUserCourse deletes a UserCourse
func (r *BaseRepository) DeleteUserCourse(id int64) error {
	return r.client.Delete(r.ctx, datastore.IDKey("UserCourse", id, nil))
}

func (r *BaseRepository) UpdateUserCourse(id int64, UserCourse *models.UserCourse) (*datastore.Key, error) {
	return r.client.Put(r.ctx, datastore.IDKey("UserCourse", id, nil), UserCourse)
}

func (r *BaseRepository) GetUserCoursesByUserID(userID int64) ([]*models.UserCourse, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("UserID", "=", userID)
	_, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	return UserCourses, nil
}

func (r *BaseRepository) GetUserCoursesByCourseID(courseID int64) ([]*models.UserCourse, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("CourseID", "=", courseID)
	_, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	return UserCourses, nil
}

func (r *BaseRepository) GetUserCourseByUserIDAndCourseID(userID int64, courseID int64) (*models.UserCourse, error) {

	// TODO: currently using an array but this should only return one
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("UserID", "=", userID).FilterField("CourseID", "=", courseID)
	keys, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	// add key values to UserCourses
	for i, key := range keys {
		UserCourses[i].KeyID = key.ID
	}

	return UserCourses[0], nil
}

func (r *BaseRepository) GetCoursesByUserID(userID int64) ([]*models.Course, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("UserID", "=", userID)
	_, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	var Courses []*models.Course
	for _, UserCourse := range UserCourses {
		Course, err := r.GetCourseByID(UserCourse.CourseID)
		if err != nil {
			return nil, err
		}

		Course.KeyID = UserCourse.CourseID
		Courses = append(Courses, Course)
	}

	return Courses, nil
}

func (r *BaseRepository) GetUsersByCourseID(courseID int64) ([]*models.User, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("CourseID", "=", courseID).FilterField("Role", "!=", "instructor")
	_, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	var Users []*models.User
	for _, UserCourse := range UserCourses {
		User, err := r.GetUserByID(UserCourse.UserID)
		if err != nil {
			return nil, err
		}
		Users = append(Users, User)
	}

	return Users, nil
}

func (r *BaseRepository) GetInstructorsByCourseID(courseID int64) ([]*models.User, error) {
	var UserCourses []*models.UserCourse
	query := datastore.NewQuery("UserCourse").FilterField("CourseID", "=", courseID).FilterField("Role", "=", "instructor")
	_, err := r.client.GetAll(r.ctx, query, &UserCourses)
	if err != nil {
		return nil, err
	}

	var Users []*models.User
	for _, UserCourse := range UserCourses {
		User, err := r.GetUserByID(UserCourse.UserID)
		if err != nil {
			return nil, err
		}
		Users = append(Users, User)
	}

	return Users, nil
}
