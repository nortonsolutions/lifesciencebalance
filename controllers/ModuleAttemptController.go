package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ModuleAttemptHandler manages user attempts and submissions for modules
type ModuleAttemptHandler struct {
	moduleRepo        models.ModuleRepository
	elementRepo       models.ElementRepository
	moduleElementRepo models.ModuleElementRepository
	userRepo          models.UserRepository
}

// ModuleSubmission represents a user's module submission with answers
type ModuleSubmission struct {
	ModuleID  int64                    `json:"module_id"`
	Answers   map[string]models.Answer `json:"answers"`
	TimeSpent int                      `json:"time_spent"`
}

// ModuleResult represents the result of a module submission
type ModuleResult struct {
	Score        int    `json:"score"`
	MaxScore     int    `json:"max_score"`
	PassingScore int    `json:"passing_score"`
	Passed       bool   `json:"passed"`
	Feedback     string `json:"feedback"`
}

// NewModuleAttemptHandler creates a new module attempt handler
func NewModuleAttemptHandler(moduleRepo models.ModuleRepository, elementRepo models.ElementRepository,
	moduleElementRepo models.ModuleElementRepository, userRepo models.UserRepository) *ModuleAttemptHandler {
	return &ModuleAttemptHandler{
		moduleRepo:        moduleRepo,
		elementRepo:       elementRepo,
		moduleElementRepo: moduleElementRepo,
		userRepo:          userRepo,
	}
}

// StartModule initializes a module session for a user
func (h *ModuleAttemptHandler) StartModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Get the module to get time limit and other settings
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Get the module elements (questions/content)
	moduleElements, err := h.moduleElementRepo.GetModuleElementsByModuleID(moduleID)
	if err != nil {
		http.Error(w, "Failed to retrieve module elements", http.StatusInternalServerError)
		return
	}

	// Sort module elements by SortKey
	sortedModuleElements := make([]*models.ModuleElement, len(moduleElements))
	for i, me := range moduleElements {
		sortedModuleElements[i] = me
	}

	// Simple bubble sort by SortKey (could be more efficient but typically module elements are few)
	for i := 0; i < len(sortedModuleElements); i++ {
		for j := 0; j < len(sortedModuleElements)-i-1; j++ {
			if sortedModuleElements[j].SortKey > sortedModuleElements[j+1].SortKey {
				sortedModuleElements[j], sortedModuleElements[j+1] = sortedModuleElements[j+1], sortedModuleElements[j]
			}
		}
	}

	// Get the actual elements
	elements := make([]*models.Element, 0)
	for _, me := range sortedModuleElements {
		element, err := h.elementRepo.GetElementByID(me.ElementID)
		if err == nil {
			// Remove correct answers from choices before sending to client
			for i := range element.Choices {
				element.Choices[i].Correct = false
			}
			elements = append(elements, element)
		}
	}

	// Create a module session response
	moduleSession := struct {
		Module    *models.Module    `json:"module"`
		Elements  []*models.Element `json:"elements"`
		StartTime time.Time         `json:"start_time"`
	}{
		Module:    module,
		Elements:  elements,
		StartTime: time.Now(),
	}

	// Return the module session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleSession)
}

// SubmitModule handles a module submission and returns the results
func (h *ModuleAttemptHandler) SubmitModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decode the submission
	var submission ModuleSubmission
	err = json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid module submission", http.StatusBadRequest)
		return
	}

	// Get the module for grading criteria
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Get the user to update their modules
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get all elements for this module for grading
	moduleElements, err := h.moduleElementRepo.GetModuleElementsByModuleID(moduleID)
	if err != nil {
		http.Error(w, "Failed to retrieve module elements", http.StatusInternalServerError)
		return
	}

	// Grade the module submission
	score, maxScore := 0, len(moduleElements)
	for _, me := range moduleElements {
		element, err := h.elementRepo.GetElementByID(me.ElementID)
		if err != nil {
			continue
		}

		// Skip non-question elements (e.g. content only elements)
		if element.Type == "content" || element.Type == "" {
			maxScore--
			continue
		}

		// Look up the user's answer
		elementIDStr := strconv.FormatInt(me.ElementID, 10)
		userAnswer, found := submission.Answers[elementIDStr]
		if !found {
			continue // Skip if no answer provided
		}

		// Grade based on question type
		switch element.Type {
		case "single", "multiple":
			correct := true
			for i, choice := range element.Choices {
				// For single/multiple choice, check if user's selected choices match correct ones
				if i < len(userAnswer.Answer) && userAnswer.Answer[i] != choice.Correct {
					correct = false
					break
				}
			}
			if correct {
				score++
				userAnswer.Correct = true
				submission.Answers[elementIDStr] = userAnswer
			}
		case "text":
			// For text questions, check against regex if implemented
			if element.TextRegex != "" {
				// Here we would implement regex matching
				// For now, just mark as correct if there's an answer
				if userAnswer.AnswerText != "" {
					score++
					userAnswer.Correct = true
					submission.Answers[elementIDStr] = userAnswer
				}
			}
		case "essay":
			// Essays would need manual grading, mark as pending
			// For now, just count it as a point if there's any content
			if userAnswer.AnswerEssay != "" {
				score++
				userAnswer.Correct = true
				submission.Answers[elementIDStr] = userAnswer
			}
		case "project":
			// Project would need manual grading or separate review
			if userAnswer.ProjectID > 0 {
				score++
				userAnswer.Correct = true
				submission.Answers[elementIDStr] = userAnswer
			}
		}
	}

	// Calculate percentage
	var percentage int
	if maxScore > 0 {
		percentage = (score * 100) / maxScore
	}

	// Determine if the user passed
	passed := percentage >= module.MinPassing

	// Create user module record
	userModule := models.UserModule{
		UserID:     userID,
		ModuleID:   moduleID,
		Answers:    submission.Answers,
		Date:       time.Now().Format(time.RFC3339),
		Score:      percentage,
		TimePassed: submission.TimeSpent,
	}

	// Update user's modules
	if user.Modules == nil {
		user.Modules = []models.UserModule{}
	}

	// Check if this module already exists for this user
	moduleFound := false
	for i, m := range user.Modules {
		if m.ModuleID == moduleID {
			user.Modules[i] = userModule
			moduleFound = true
			break
		}
	}

	if !moduleFound {
		user.Modules = append(user.Modules, userModule)
	}

	// Update the user
	_, err = h.userRepo.UpdateUser(userID, user)
	if err != nil {
		http.Error(w, "Failed to update user progress", http.StatusInternalServerError)
		return
	}

	// Prepare the result
	result := ModuleResult{
		Score:        score,
		MaxScore:     maxScore,
		PassingScore: module.MinPassing,
		Passed:       passed,
	}

	if passed {
		result.Feedback = "Congratulations! You completed this module successfully."
	} else {
		result.Feedback = "You did not meet the passing criteria for this module. Please review the material and try again."
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetModuleResults gets a user's results for a specific module
func (h *ModuleAttemptHandler) GetModuleResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the user to check their modules
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Find the module result
	var userModule *models.UserModule
	for i, m := range user.Modules {
		if m.ModuleID == moduleID {
			userModule = &user.Modules[i]
			break
		}
	}

	if userModule == nil {
		http.Error(w, "No results found for this module", http.StatusNotFound)
		return
	}

	// Get the module for reference
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Get all elements for this module
	elements, err := h.moduleElementRepo.GetElementsByModuleID(moduleID)
	if err != nil {
		http.Error(w, "Failed to retrieve module elements", http.StatusInternalServerError)
		return
	}

	// Return the results
	result := struct {
		Module     *models.Module     `json:"module"`
		UserModule *models.UserModule `json:"user_module"`
		Elements   []*models.Element  `json:"elements"`
		Passed     bool               `json:"passed"`
	}{
		Module:     module,
		UserModule: userModule,
		Elements:   elements,
		Passed:     userModule.Score >= module.MinPassing,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetModuleAnalytics gets aggregate statistics for a module
func (h *ModuleAttemptHandler) GetModuleAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Get all users
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Get the module
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Analyze module attempts
	totalAttempts := 0
	totalPassed := 0
	averageScore := 0
	averageTime := 0

	// Track element-specific statistics
	elementStats := make(map[string]struct {
		Attempts       int
		Correct        int
		PercentCorrect float64
	})

	// Collect data
	for _, user := range users {
		for _, m := range user.Modules {
			if m.ModuleID == moduleID {
				totalAttempts++
				if m.Score >= module.MinPassing {
					totalPassed++
				}
				averageScore += m.Score
				averageTime += m.TimePassed

				// Element-specific stats
				for elementID, answer := range m.Answers {
					stats := elementStats[elementID]
					stats.Attempts++
					if answer.Correct {
						stats.Correct++
					}
					if stats.Attempts > 0 {
						stats.PercentCorrect = float64(stats.Correct) / float64(stats.Attempts) * 100
					}
					elementStats[elementID] = stats
				}
			}
		}
	}

	// Calculate averages
	if totalAttempts > 0 {
		averageScore = averageScore / totalAttempts
		averageTime = averageTime / totalAttempts
	}

	// Format the element stats to include JSON tags for the response
	formattedElementStats := make(map[string]struct {
		Attempts       int     `json:"attempts"`
		Correct        int     `json:"correct"`
		PercentCorrect float64 `json:"percent_correct"`
	})

	for id, stats := range elementStats {
		formattedElementStats[id] = struct {
			Attempts       int     `json:"attempts"`
			Correct        int     `json:"correct"`
			PercentCorrect float64 `json:"percent_correct"`
		}{
			Attempts:       stats.Attempts,
			Correct:        stats.Correct,
			PercentCorrect: stats.PercentCorrect,
		}
	}

	// Prepare the analytics result
	analytics := struct {
		Module        *models.Module `json:"module"`
		TotalAttempts int            `json:"total_attempts"`
		TotalPassed   int            `json:"total_passed"`
		PassRate      float64        `json:"pass_rate"`
		AverageScore  int            `json:"average_score"`
		AverageTime   int            `json:"average_time"` // in seconds
		ElementStats  map[string]struct {
			Attempts       int     `json:"attempts"`
			Correct        int     `json:"correct"`
			PercentCorrect float64 `json:"percent_correct"`
		} `json:"element_stats"`
	}{
		Module:        module,
		TotalAttempts: totalAttempts,
		TotalPassed:   totalPassed,
		AverageScore:  averageScore,
		AverageTime:   averageTime,
		ElementStats:  formattedElementStats,
	}

	if totalAttempts > 0 {
		analytics.PassRate = float64(totalPassed) / float64(totalAttempts) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// ResetModuleAttempt allows a user to reset their module attempt
func (h *ModuleAttemptHandler) ResetModuleAttempt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the user
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remove this module from the user's modules
	if user.Modules != nil {
		updatedModules := []models.UserModule{}
		for _, m := range user.Modules {
			if m.ModuleID != moduleID {
				updatedModules = append(updatedModules, m)
			}
		}
		user.Modules = updatedModules
	}

	// Update the user
	_, err = h.userRepo.UpdateUser(userID, user)
	if err != nil {
		http.Error(w, "Failed to reset module attempt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Module attempt reset successfully"})
}
