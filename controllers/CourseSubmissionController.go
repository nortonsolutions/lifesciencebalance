package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// CourseSubmissionHandler handles course submission and approval process
type CourseSubmissionHandler struct {
	courseRepo models.CourseRepository
	userRepo   models.UserRepository
}

// NewCourseSubmissionHandler creates a new CourseSubmissionHandler
func NewCourseSubmissionHandler(courseRepo models.CourseRepository, userRepo models.UserRepository) *CourseSubmissionHandler {
	return &CourseSubmissionHandler{
		courseRepo: courseRepo,
		userRepo:   userRepo,
	}
}

// SubmitCourse handles the submission of a new course
func (h *CourseSubmissionHandler) SubmitCourse(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the session
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Just log the userID to avoid the "unused" error
	fmt.Printf("User ID submitting course: %d\n", userID)

	// Create a new course with pending approval
	var course models.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid course data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set initial values for the course
	course.Approved = false // Not approved by default

	// Validate department
	if course.Department == "" {
		http.Error(w, "Department is required", http.StatusBadRequest)
		return
	}

	// Save the new course
	key, err := h.courseRepo.CreateCourse(&course)
	if err != nil {
		http.Error(w, "Failed to create course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	course.KeyID = key.ID

	// Return the created course
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

// GetPendingCourses returns all pending (unapproved) courses
func (h *CourseSubmissionHandler) GetPendingCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courseRepo.GetUnapprovedCourses()
	if err != nil {
		http.Error(w, "Failed to get pending courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// ApproveCourse approves a course submission
func (h *CourseSubmissionHandler) ApproveCourse(w http.ResponseWriter, r *http.Request) {
	// Get the course ID from the URL
	vars := mux.Vars(r)
	courseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Get the course
	course, err := h.courseRepo.GetCourseByID(courseID)
	if err != nil {
		http.Error(w, "Course not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Update the course status
	course.Approved = true

	// Save the updated course
	_, err = h.courseRepo.UpdateCourse(courseID, course)
	if err != nil {
		http.Error(w, "Failed to approve course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "approved"})
}

// RejectCourse rejects a course submission
func (h *CourseSubmissionHandler) RejectCourse(w http.ResponseWriter, r *http.Request) {
	// Get the course ID from the URL
	vars := mux.Vars(r)
	courseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Get rejection reason from request body
	var requestBody struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// If there's an error, proceed without a reason
		requestBody.Reason = ""
	}

	// Get the course
	course, err := h.courseRepo.GetCourseByID(courseID)
	if err != nil {
		http.Error(w, "Course not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Update the course status
	course.Approved = false
	// Store rejection reason if provided
	if requestBody.Reason != "" {
		// You might want to add a RejectionReason field to your Course model
		// For now, we'll just demonstrate the concept
		fmt.Printf("Rejection reason: %s\n", requestBody.Reason)
	}

	// Save the updated course
	_, err = h.courseRepo.UpdateCourse(courseID, course)
	if err != nil {
		http.Error(w, "Failed to reject course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "rejected"})
}

// GetApprovedCourses returns all approved courses
func (h *CourseSubmissionHandler) GetApprovedCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courseRepo.GetApprovedCourses()
	if err != nil {
		http.Error(w, "Failed to get approved courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCoursesByDepartment returns courses filtered by department
func (h *CourseSubmissionHandler) GetCoursesByDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	department := vars["department"]

	courses, err := h.courseRepo.GetCoursesByDepartment(department)
	if err != nil {
		http.Error(w, "Failed to get courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
