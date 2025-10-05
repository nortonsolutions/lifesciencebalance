package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// DepartmentHandler handles department operations
type DepartmentHandler struct {
	// In a real application, this would use a proper repository
	// For this example, we'll keep it simple with in-memory storage
	departments []string
}

// NewDepartmentHandler creates a new DepartmentHandler
func NewDepartmentHandler() *DepartmentHandler {
	// Initialize with default departments
	return &DepartmentHandler{
		departments: []string{
			"Medicine",
			"Biology",
			"Chemistry",
			"Nursing",
			"Public Health",
			"Research Methodology",
			"Healthcare Management",
		},
	}
}

// GetAllDepartments returns all departments
func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	type Department struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// Convert string array to Department objects
	departmentObjects := make([]Department, len(h.departments))
	for i, dept := range h.departments {
		departmentObjects[i] = Department{
			ID:   i,
			Name: dept,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departmentObjects)
}

// AddDepartment adds a new department
func (h *DepartmentHandler) AddDepartment(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate department name
	departmentName := strings.TrimSpace(requestBody.Name)
	if departmentName == "" {
		http.Error(w, "Department name cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if department already exists
	for _, dept := range h.departments {
		if strings.ToLower(dept) == strings.ToLower(departmentName) {
			http.Error(w, "Department already exists", http.StatusConflict)
			return
		}
	}

	// Add department
	h.departments = append(h.departments, departmentName)

	// Create response
	response := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		ID:   len(h.departments) - 1,
		Name: departmentName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteDepartment deletes a department
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would get the ID from the URL path
	// For this example, we'll parse it from the query string
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}

	// Convert ID to int
	var departmentID int
	if _, err := fmt.Sscanf(id, "%d", &departmentID); err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	// Validate ID
	if departmentID < 0 || departmentID >= len(h.departments) {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	// Remove department (this is a simple implementation)
	h.departments = append(h.departments[:departmentID], h.departments[departmentID+1:]...)

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
