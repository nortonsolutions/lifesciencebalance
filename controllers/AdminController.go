package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// AdminHandler handles administrative operations
type AdminHandler struct {
	userRepository    models.UserRepository
	courseRepository  models.CourseRepository
	moduleRepository  models.ModuleRepository
	elementRepository models.ElementRepository
	projectRepository models.ProjectRepository
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(userRepository models.UserRepository, courseRepository models.CourseRepository,
	moduleRepository models.ModuleRepository, elementRepository models.ElementRepository,
	projectRepository models.ProjectRepository) *AdminHandler {
	return &AdminHandler{
		userRepository:    userRepository,
		courseRepository:  courseRepository,
		moduleRepository:  moduleRepository,
		elementRepository: elementRepository,
		projectRepository: projectRepository,
	}
}

// GetSystemStats retrieves high-level system statistics
func (h *AdminHandler) GetSystemStats(w http.ResponseWriter, r *http.Request) {
	// Get counts of various entities
	users, err := h.userRepository.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	courses, err := h.courseRepository.GetAllCourses()
	if err != nil {
		http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
		return
	}

	modules, err := h.moduleRepository.GetAllModules()
	if err != nil {
		http.Error(w, "Failed to retrieve modules", http.StatusInternalServerError)
		return
	}

	elements, err := h.elementRepository.GetAllElements()
	if err != nil {
		http.Error(w, "Failed to retrieve elements", http.StatusInternalServerError)
		return
	}

	projects, err := h.projectRepository.GetAllProjects()
	if err != nil {
		http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
		return
	}

	// Count active users (users who have completed at least one module)
	activeUsers := 0
	for _, user := range users {
		if len(user.Modules) > 0 {
			activeUsers++
		}
	}

	// Count completed modules
	completedModules := 0
	for _, user := range users {
		completedModules += len(user.Modules)
	}

	// Count modules by purpose (based on properties)
	quizModules := 0
	projectModules := 0
	contentModules := 0

	for _, module := range modules {
		// Classify modules based on their properties
		if module.TimeLimit > 0 || module.MaxAttempts > 0 {
			quizModules++ // If it has a time limit or max attempts, it's likely a quiz
		} else {
			// Look for project submissions from this module
			projects, _ := h.projectRepository.GetProjectsByModuleID(module.KeyID)
			if len(projects) > 0 {
				projectModules++
			} else {
				contentModules++ // Default to content module
			}
		}
	}

	// Prepare stats
	stats := map[string]interface{}{
		"users": map[string]interface{}{
			"total":  len(users),
			"active": activeUsers,
		},
		"courses": map[string]interface{}{
			"total": len(courses),
		},
		"modules": map[string]interface{}{
			"total":     len(modules),
			"quiz":      quizModules,
			"project":   projectModules,
			"content":   contentModules,
			"completed": completedModules,
		},
		"elements": map[string]interface{}{
			"total": len(elements),
		},
		"projects": map[string]interface{}{
			"total": len(projects),
		},
	}

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetUserStats retrieves statistics for a specific user
func (h *AdminHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
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

	// Count completed modules
	completedModules := len(user.Modules)

	// Calculate average score
	totalScore := 0
	for _, module := range user.Modules {
		totalScore += module.Score
	}

	averageScore := 0
	if completedModules > 0 {
		averageScore = totalScore / completedModules
	}

	// Count projects submitted
	projects, err := h.projectRepository.GetProjectsByUserID(userID)
	if err != nil {
		projects = []*models.Project{}
	}

	// Prepare stats
	stats := map[string]interface{}{
		"user": map[string]interface{}{
			"id":        user.KeyID,
			"username":  user.Username,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"roles":     user.Roles,
		},
		"progress": map[string]interface{}{
			"completed_modules":  completedModules,
			"average_score":      averageScore,
			"projects_submitted": len(projects),
		},
	}

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetCourseStats retrieves statistics for a specific course
func (h *AdminHandler) GetCourseStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Get the course
	course, err := h.courseRepository.GetCourseByID(courseID)
	if err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Get modules for this course
	modules, err := h.moduleRepository.GetAllModulesByCourseID(courseID)
	if err != nil {
		modules = []*models.Module{}
	}

	// Get all users
	users, err := h.userRepository.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Count users enrolled in this course
	enrolledUsers := 0
	completedUsers := 0
	moduleCompletions := make(map[int64]int)

	for _, user := range users {
		userEnrolled := false
		userCompletedAll := true

		for _, module := range user.Modules {
			// Check if module belongs to this course
			for _, courseModule := range modules {
				if module.ModuleID == courseModule.KeyID {
					userEnrolled = true
					moduleCompletions[module.ModuleID]++

					// Check if user passed the module
					if module.Score < courseModule.MinPassing {
						userCompletedAll = false
					}
				}
			}
		}

		if userEnrolled {
			enrolledUsers++
			if userCompletedAll && len(modules) > 0 {
				completedUsers++
			}
		}
	}

	// Prepare stats
	stats := map[string]interface{}{
		"course": map[string]interface{}{
			"id":   course.KeyID,
			"name": course.Name,
		},
		"modules": map[string]interface{}{
			"total": len(modules),
		},
		"users": map[string]interface{}{
			"enrolled":  enrolledUsers,
			"completed": completedUsers,
		},
		"module_completions": moduleCompletions,
	}

	// Return stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
