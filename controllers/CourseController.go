package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// CourseHandler ..
type CourseHandler struct {
	courseRepository models.CourseRepository
}

// NewCourseHandler ..
func NewCourseHandler(courseRepository models.CourseRepository) *CourseHandler {
	return &CourseHandler{courseRepository: courseRepository}
}

// add course
func (c *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	course := models.Course{}
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.courseRepository.CreateCourse(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key.ID)
}

// remove course
func (c *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.courseRepository.DeleteCourse(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted"})
}

// get all courses
func (c *CourseHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := c.courseRepository.GetAllCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// get course by id
func (c *CourseHandler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	course, err := c.courseRepository.GetCourseByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// update course
func (c *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	course := models.Course{}
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.courseRepository.UpdateCourse(idInt, &course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// GetApprovedCourses returns only approved courses
func (c *CourseHandler) GetApprovedCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := c.courseRepository.GetApprovedCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetUnapprovedCourses returns only unapproved courses
func (c *CourseHandler) GetUnapprovedCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := c.courseRepository.GetUnapprovedCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCoursesByDepartment returns courses filtered by department
func (c *CourseHandler) GetCoursesByDepartment(w http.ResponseWriter, r *http.Request) {
	var department = mux.Vars(r)["department"]
	courses, err := c.courseRepository.GetCoursesByDepartment(department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// ApproveCourse approves a course
func (c *CourseHandler) ApproveCourse(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	// Get the course first
	course, err := c.courseRepository.GetCourseByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set approved to true
	course.Approved = true

	// Update the course
	key, err := c.courseRepository.UpdateCourse(idInt, course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// UnapproveeCourse unapproves a course
func (c *CourseHandler) UnapproveeCourse(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	// Get the course first
	course, err := c.courseRepository.GetCourseByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set approved to false
	course.Approved = false

	// Update the course
	key, err := c.courseRepository.UpdateCourse(idInt, course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}
