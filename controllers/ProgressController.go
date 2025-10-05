package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ProgressHandler ..
type ProgressHandler struct {
	userRepository       models.UserRepository
	courseRepository     models.CourseRepository
	moduleRepository     models.ModuleRepository
	userCourseRepository models.UserCourseRepository
}

// NewProgressHandler ..
func NewProgressHandler(userRepository models.UserRepository, courseRepository models.CourseRepository,
	moduleRepository models.ModuleRepository, userCourseRepository models.UserCourseRepository) *ProgressHandler {
	return &ProgressHandler{
		userRepository:       userRepository,
		courseRepository:     courseRepository,
		moduleRepository:     moduleRepository,
		userCourseRepository: userCourseRepository,
	}
}

// UserProgressSummary represents a summary of the user's progress in all courses
type UserProgressSummary struct {
	TotalCourses      int                     `json:"total_courses"`
	CompletedCourses  int                     `json:"completed_courses"`
	InProgressCourses int                     `json:"in_progress_courses"`
	NotStartedCourses int                     `json:"not_started_courses"`
	CoursesProgress   []CourseProgressSummary `json:"courses_progress"`
	OverallProgress   int                     `json:"overall_progress"` // percentage
	CompletedQuizzes  int                     `json:"completed_quizzes"`
	TotalQuizzes      int                     `json:"total_quizzes"`
	PassedQuizzes     int                     `json:"passed_quizzes"`
}

// CourseProgressSummary represents a summary of progress in a specific course
type CourseProgressSummary struct {
	CourseID         int64  `json:"course_id"`
	CourseName       string `json:"course_name"`
	ModulesCompleted int    `json:"modules_completed"`
	TotalModules     int    `json:"total_modules"`
	Progress         int    `json:"progress"` // percentage
	StartDate        string `json:"start_date,omitempty"`
	CompletionDate   string `json:"completion_date,omitempty"`
	Grade            int    `json:"grade,omitempty"`
}

// GetUserProgress returns the user's overall progress
func (h *ProgressHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the user
	user, err := h.userRepository.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get the user's courses
	userCourses, err := h.userCourseRepository.GetUserCoursesByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user courses", http.StatusInternalServerError)
		return
	}

	// Create a progress summary
	summary := UserProgressSummary{
		TotalCourses:      len(userCourses),
		CompletedCourses:  0,
		InProgressCourses: 0,
		NotStartedCourses: 0,
		CoursesProgress:   make([]CourseProgressSummary, 0),
	}

	totalQuizzes := 0
	completedQuizzes := 0
	passedQuizzes := 0

	// Process each course
	for _, userCourse := range userCourses {
		course, err := h.courseRepository.GetCourseByID(userCourse.CourseID)
		if err != nil {
			continue
		}

		// Get all modules for this course
		modules, err := h.moduleRepository.GetAllModulesByCourseID(course.KeyID)
		if err != nil {
			continue
		}

		totalModules := len(modules)
		if totalModules == 0 {
			continue
		}

		// Count completed modules
		completedModules := 0
		if user.Modules != nil {
			for _, userModule := range user.Modules {
				for _, module := range modules {
					if userModule.ModuleID == module.KeyID && userModule.Score >= module.MinPassing {
						completedModules++
						completedQuizzes++
						if userModule.Score >= module.MinPassing {
							passedQuizzes++
						}
					}
				}
			}
		}

		totalQuizzes += totalModules

		// Calculate course progress
		courseProgress := 0
		if totalModules > 0 {
			courseProgress = (completedModules * 100) / totalModules
		}

		// Determine course status
		if userCourse.StartedOn.IsZero() {
			summary.NotStartedCourses++
		} else if !userCourse.CompletedOn.IsZero() {
			summary.CompletedCourses++
		} else {
			summary.InProgressCourses++
		}

		// Add course summary
		courseSummary := CourseProgressSummary{
			CourseID:         course.KeyID,
			CourseName:       course.Name,
			ModulesCompleted: completedModules,
			TotalModules:     totalModules,
			Progress:         courseProgress,
			Grade:            userCourse.Grade,
		}

		if !userCourse.StartedOn.IsZero() {
			courseSummary.StartDate = userCourse.StartedOn.Format(time.RFC3339)
		}
		if !userCourse.CompletedOn.IsZero() {
			courseSummary.CompletionDate = userCourse.CompletedOn.Format(time.RFC3339)
		}

		summary.CoursesProgress = append(summary.CoursesProgress, courseSummary)
	}

	// Calculate overall progress
	if summary.TotalCourses > 0 {
		totalProgress := 0
		for _, course := range summary.CoursesProgress {
			totalProgress += course.Progress
		}
		summary.OverallProgress = totalProgress / summary.TotalCourses
	}

	summary.CompletedQuizzes = completedQuizzes
	summary.TotalQuizzes = totalQuizzes
	summary.PassedQuizzes = passedQuizzes

	// Return the summary
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetCourseProgress returns the user's progress in a specific course
func (h *ProgressHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseInt(vars["courseId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Get the user
	user, err := h.userRepository.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get the user course record
	userCourse, err := h.userCourseRepository.GetUserCourseByUserIDAndCourseID(userID, courseID)
	if err != nil {
		http.Error(w, "User is not enrolled in this course", http.StatusNotFound)
		return
	}

	// Get the course
	course, err := h.courseRepository.GetCourseByID(courseID)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Get all modules for this course
	modules, err := h.moduleRepository.GetAllModulesByCourseID(courseID)
	if err != nil {
		http.Error(w, "Failed to retrieve course modules", http.StatusInternalServerError)
		return
	}

	// Create a detailed course progress report
	type ModuleProgress struct {
		ModuleID      int64  `json:"module_id"`
		ModuleName    string `json:"module_name"`
		Completed     bool   `json:"completed"`
		Score         int    `json:"score,omitempty"`
		AttemptCount  int    `json:"attempt_count,omitempty"`
		DateCompleted string `json:"date_completed,omitempty"`
	}

	courseProgress := struct {
		CourseID         int64            `json:"course_id"`
		CourseName       string           `json:"course_name"`
		ModulesCompleted int              `json:"modules_completed"`
		TotalModules     int              `json:"total_modules"`
		Progress         int              `json:"progress"` // percentage
		Grade            int              `json:"grade,omitempty"`
		StartDate        string           `json:"start_date,omitempty"`
		CompletionDate   string           `json:"completion_date,omitempty"`
		ModuleProgress   []ModuleProgress `json:"module_progress"`
	}{
		CourseID:     course.KeyID,
		CourseName:   course.Name,
		TotalModules: len(modules),
		Grade:        userCourse.Grade,
	}

	if !userCourse.StartedOn.IsZero() {
		courseProgress.StartDate = userCourse.StartedOn.Format(time.RFC3339)
	}
	if !userCourse.CompletedOn.IsZero() {
		courseProgress.CompletionDate = userCourse.CompletedOn.Format(time.RFC3339)
	}

	// Process each module
	completedModules := 0
	courseProgress.ModuleProgress = make([]ModuleProgress, len(modules))

	for i, module := range modules {
		moduleProgress := ModuleProgress{
			ModuleID:   module.KeyID,
			ModuleName: module.Name,
			Completed:  false,
		}

		// Check if user has completed this module
		if user.Modules != nil {
			for _, userModule := range user.Modules {
				if userModule.ModuleID == module.KeyID {
					moduleProgress.Score = userModule.Score
					moduleProgress.AttemptCount = 1 // Just a placeholder, we don't track attempts yet
					moduleProgress.Completed = userModule.Score >= module.MinPassing
					if moduleProgress.Completed {
						completedModules++
						moduleProgress.DateCompleted = userModule.Date
					}
					break
				}
			}
		}

		courseProgress.ModuleProgress[i] = moduleProgress
	}

	courseProgress.ModulesCompleted = completedModules
	if courseProgress.TotalModules > 0 {
		courseProgress.Progress = (completedModules * 100) / courseProgress.TotalModules
	}

	// Return the course progress
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courseProgress)
}

// UpdateUserCourseProgress updates the progress of a user in a course
func (h *ProgressHandler) UpdateUserCourseProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseInt(vars["courseId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Get the user course record
	userCourse, err := h.userCourseRepository.GetUserCourseByUserIDAndCourseID(userID, courseID)
	if err != nil {
		http.Error(w, "User is not enrolled in this course", http.StatusNotFound)
		return
	}

	// Get the course modules
	_, err = h.courseRepository.GetCourseByID(courseID)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	modules, err := h.moduleRepository.GetAllModulesByCourseID(courseID)
	if err != nil {
		http.Error(w, "Failed to retrieve course modules", http.StatusInternalServerError)
		return
	}

	if len(modules) == 0 {
		http.Error(w, "Course has no modules", http.StatusBadRequest)
		return
	}

	// Get the user
	user, err := h.userRepository.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the user has completed all modules
	totalModules := len(modules)
	completedModules := 0
	totalScore := 0

	if user.Modules != nil {
		for _, userModule := range user.Modules {
			for _, module := range modules {
				if userModule.ModuleID == module.KeyID && userModule.Score >= module.MinPassing {
					completedModules++
					totalScore += userModule.Score
					break
				}
			}
		}
	}

	// Update course progress
	if completedModules > 0 {
		// If not started yet, set start date
		if userCourse.StartedOn.IsZero() {
			userCourse.StartedOn = time.Now()
		}

		// Calculate overall grade
		userCourse.Grade = totalScore / completedModules

		// If all modules completed, mark course as completed
		if completedModules == totalModules {
			userCourse.CompletedOn = time.Now()
		}

		// Update the user course record
		_, err = h.userCourseRepository.UpdateUserCourse(userCourse.KeyID, userCourse)
		if err != nil {
			http.Error(w, "Failed to update course progress", http.StatusInternalServerError)
			return
		}
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Course progress updated",
		"progress": map[string]interface{}{
			"completed_modules": completedModules,
			"total_modules":     totalModules,
			"progress":          float64(completedModules) / float64(totalModules) * 100,
			"grade":             userCourse.Grade,
		},
	})
}
