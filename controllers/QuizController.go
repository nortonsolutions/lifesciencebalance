package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type QuizHandler struct {
	moduleRepo        models.ModuleRepository
	elementRepo       models.ElementRepository
	moduleElementRepo models.ModuleElementRepository
	userRepo          models.UserRepository
}

// QuizSubmission represents a user's quiz submission
type QuizSubmission struct {
	ModuleID  int64                    `json:"module_id"`
	Answers   map[string]models.Answer `json:"answers"`
	TimeSpent int                      `json:"time_spent"`
}

// QuizResult represents the result of a quiz submission
type QuizResult struct {
	Score        int    `json:"score"`
	MaxScore     int    `json:"max_score"`
	PassingScore int    `json:"passing_score"`
	Passed       bool   `json:"passed"`
	Feedback     string `json:"feedback"`
}

func NewQuizHandler(moduleRepo models.ModuleRepository, elementRepo models.ElementRepository,
	moduleElementRepo models.ModuleElementRepository, userRepo models.UserRepository) *QuizHandler {
	return &QuizHandler{
		moduleRepo:        moduleRepo,
		elementRepo:       elementRepo,
		moduleElementRepo: moduleElementRepo,
		userRepo:          userRepo,
	}
}

// StartQuiz initializes a quiz session for a user
func (h *QuizHandler) StartQuiz(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Get the module to get time limit and other quiz settings
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Get the module elements (questions)
	moduleElements, err := h.moduleElementRepo.GetModuleElementsByModuleID(moduleID)
	if err != nil {
		http.Error(w, "Failed to retrieve quiz questions", http.StatusInternalServerError)
		return
	}

	// Get the actual elements
	elementIDs := make([]int64, len(moduleElements))
	for i, me := range moduleElements {
		elementIDs[i] = me.ElementID
	}

	elements := make([]*models.Element, 0)
	for _, elementID := range elementIDs {
		element, err := h.elementRepo.GetElementByID(elementID)
		if err == nil {
			// Remove correct answers from choices before sending to client
			for i := range element.Choices {
				element.Choices[i].Correct = false
			}
			elements = append(elements, element)
		}
	}

	// Create a quiz session response
	quizSession := struct {
		Module    *models.Module    `json:"module"`
		Elements  []*models.Element `json:"elements"`
		StartTime time.Time         `json:"start_time"`
	}{
		Module:    module,
		Elements:  elements,
		StartTime: time.Now(),
	}

	// Return the quiz session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizSession)
}

// SubmitQuiz handles a quiz submission and returns the results
func (h *QuizHandler) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
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
	var submission QuizSubmission
	err = json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid quiz submission", http.StatusBadRequest)
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
		http.Error(w, "Failed to retrieve quiz elements", http.StatusInternalServerError)
		return
	}

	// Grade the quiz
	score, maxScore := 0, len(moduleElements)
	for _, me := range moduleElements {
		element, err := h.elementRepo.GetElementByID(me.ElementID)
		if err != nil {
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
				// For now, just mark as correct
				score++
				userAnswer.Correct = true
				submission.Answers[elementIDStr] = userAnswer
			}
		case "essay":
			// Essays would need manual grading, mark as pending
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
	result := QuizResult{
		Score:        score,
		MaxScore:     maxScore,
		PassingScore: module.MinPassing,
		Passed:       passed,
	}

	if passed {
		result.Feedback = "Congratulations! You passed the quiz."
	} else {
		result.Feedback = "You did not pass the quiz. Please review the material and try again."
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetQuizResults gets a user's quiz results for a module
func (h *QuizHandler) GetQuizResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
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
	for _, m := range user.Modules {
		if m.ModuleID == moduleID {
			userModule = &m
			break
		}
	}

	if userModule == nil {
		http.Error(w, "No quiz results found for this module", http.StatusNotFound)
		return
	}

	// Get the module for reference
	module, err := h.moduleRepo.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Return the results
	result := struct {
		Module     *models.Module     `json:"module"`
		UserModule *models.UserModule `json:"user_module"`
		Passed     bool               `json:"passed"`
	}{
		Module:     module,
		UserModule: userModule,
		Passed:     userModule.Score >= module.MinPassing,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetQuizAnalytics gets aggregate statistics for a quiz
func (h *QuizHandler) GetQuizAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
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

	// Analyze quiz attempts
	totalAttempts := 0
	totalPassed := 0
	averageScore := 0
	averageTime := 0

	// Track question-specific statistics
	questionStats := make(map[string]struct {
		Attempts       int
		Correct        int
		PercentCorrect float64
	})

	// Collect data
	for _, user := range users {
		for _, m := range user.Modules {
			if m.ModuleID == moduleID {
				totalAttempts++
				if m.Score >= 70 { // Assuming 70% is passing
					totalPassed++
				}
				averageScore += m.Score
				averageTime += m.TimePassed

				// Question-specific stats
				for questionID, answer := range m.Answers {
					stats := questionStats[questionID]
					stats.Attempts++
					if answer.Correct {
						stats.Correct++
					}
					stats.PercentCorrect = float64(stats.Correct) / float64(stats.Attempts) * 100
					questionStats[questionID] = stats
				}
			}
		}
	}

	// Calculate averages
	if totalAttempts > 0 {
		averageScore /= totalAttempts
		averageTime /= totalAttempts
	}

	// Convert the questionStats to the format required by the struct
	formattedQuestionStats := make(map[string]struct {
		Attempts       int     `json:"attempts"`
		Correct        int     `json:"correct"`
		PercentCorrect float64 `json:"percent_correct"`
	})

	for id, stats := range questionStats {
		formattedQuestionStats[id] = struct {
			Attempts       int     `json:"attempts"`
			Correct        int     `json:"correct"`
			PercentCorrect float64 `json:"percent_correct"`
		}{
			Attempts:       stats.Attempts,
			Correct:        stats.Correct,
			PercentCorrect: stats.PercentCorrect,
		}
	}

	// Prepare the result
	analytics := struct {
		TotalAttempts int     `json:"total_attempts"`
		TotalPassed   int     `json:"total_passed"`
		PassRate      float64 `json:"pass_rate"`
		AverageScore  int     `json:"average_score"`
		AverageTime   int     `json:"average_time"`
		QuestionStats map[string]struct {
			Attempts       int     `json:"attempts"`
			Correct        int     `json:"correct"`
			PercentCorrect float64 `json:"percent_correct"`
		} `json:"question_stats"`
	}{
		TotalAttempts: totalAttempts,
		TotalPassed:   totalPassed,
		AverageScore:  averageScore,
		AverageTime:   averageTime,
		QuestionStats: formattedQuestionStats,
	}

	if totalAttempts > 0 {
		analytics.PassRate = float64(totalPassed) / float64(totalAttempts) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// ResetQuiz allows a user to reset their quiz attempt
func (h *QuizHandler) ResetQuiz(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
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
		http.Error(w, "Failed to reset quiz", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Quiz reset successfully"})
}
