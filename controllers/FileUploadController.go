package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"restAPI/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// FileUploadHandler handles file upload operations
type FileUploadHandler struct {
	projectRepository       models.ProjectRepository
	moduleRepository        models.ModuleRepository
	userRepository          models.UserRepository
	moduleElementRepository models.ModuleElementRepository
}

// FileMetadata represents metadata about an uploaded file
type FileMetadata struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	UploadDate  string `json:"upload_date"`
	FilePath    string `json:"file_path"`
}

// NewFileUploadHandler creates a new file upload handler
func NewFileUploadHandler(projectRepository models.ProjectRepository, moduleRepository models.ModuleRepository,
	userRepository models.UserRepository, moduleElementRepository models.ModuleElementRepository) *FileUploadHandler {
	return &FileUploadHandler{
		projectRepository:       projectRepository,
		moduleRepository:        moduleRepository,
		userRepository:          userRepository,
		moduleElementRepository: moduleElementRepository,
	}
}

// validateFile checks if the file is valid
func validateFile(file multipart.File, header *multipart.FileHeader) error {
	// Check file size (10MB limit)
	if header.Size > 10<<20 {
		return fmt.Errorf("file too large: maximum size is 10MB")
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".txt": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".zip": true, ".rar": true, ".7z": true,
		".mp4": true, ".mov": true, ".avi": true,
		".mp3": true, ".wav": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	return nil
}

// UploadProject handles file upload for projects
func (h *FileUploadHandler) UploadProject(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	userIDStr := r.FormValue("userId")
	moduleIDStr := r.FormValue("moduleId")
	courseIDStr := r.FormValue("courseId")
	projectName := r.FormValue("name")
	projectDescription := r.FormValue("description")

	if userIDStr == "" || moduleIDStr == "" || courseIDStr == "" || projectName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	moduleID, _ := strconv.ParseInt(moduleIDStr, 10, 64)
	courseID, _ := strconv.ParseInt(courseIDStr, 10, 64)

	// Get the uploaded file
	file, handler, err := r.FormFile("projectFile")
	if err != nil {
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate the file
	if err := validateFile(file, handler); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create upload directory if it doesn't exist
	uploadDir := "./static/uploads/projects"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a unique filename
	fileExt := filepath.Ext(handler.Filename)
	uniqueFilename := uuid.New().String() + fileExt
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Create the file on the server
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create file on server: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the file to the destination
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "Unable to save file on server: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new project record
	project := &models.Project{
		Name:        projectName,
		Description: projectDescription,
		File:        "/uploads/projects/" + uniqueFilename, // Store relative path
		Date:        time.Now(),
		UserID:      userID,
		CourseID:    courseID,
		ModuleID:    moduleID,
	}

	// Save to database
	key, err := h.projectRepository.CreateProject(project)
	if err != nil {
		http.Error(w, "Error saving project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user module record with the project reference
	user, err := h.userRepository.GetUserByID(userID)
	if err == nil && user != nil {
		// Find if user has a module record
		for i, module := range user.Modules {
			if module.ModuleID == moduleID {
				// Update existing answers if they exist
				if module.Answers == nil {
					module.Answers = make(map[string]models.Answer)
				}

				// Look for project-type elements in this module to update
				// Ideally we should know the exact element ID, but for now we'll update
				// any answer that has a ProjectID field or is empty
				updated := false

				// Get module elements to find question IDs
				elements, err := h.moduleElementRepository.GetElementsByModuleID(moduleID)
				if err == nil {
					for _, element := range elements {
						if element.Type == "project" || element.Type == "file" {
							// Found a project element, update its answer
							elementIDStr := strconv.FormatInt(element.KeyID, 10)
							answer, exists := module.Answers[elementIDStr]

							if !exists || answer.ProjectID == 0 {
								answer = models.Answer{
									ProjectID: key.ID,
									Correct:   true, // Auto-approve submission (instructor will review later)
								}
								module.Answers[elementIDStr] = answer
								updated = true
							}
						}
					}
				}

				// If no specific project element found, update first empty slot
				if !updated {
					// Create a generic element ID if none found
					elementIDStr := "project_" + strconv.FormatInt(moduleID, 10)
					module.Answers[elementIDStr] = models.Answer{
						ProjectID: key.ID,
						Correct:   true,
					}
				}

				user.Modules[i] = module

				// Update the user
				_, err = h.userRepository.UpdateUser(userID, user)
				if err != nil {
					fmt.Printf("Warning: Failed to update user with project: %v\n", err)
				}
				break
			}
		}
	}

	// Return success response with metadata
	metadata := FileMetadata{
		FileName:    handler.Filename,
		FileSize:    handler.Size,
		ContentType: handler.Header.Get("Content-Type"),
		UploadDate:  time.Now().Format(time.RFC3339),
		FilePath:    project.File,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Project uploaded successfully",
		"projectId": key.ID,
		"metadata":  metadata,
	})
}

// GetProjectFile serves the project file
func (h *FileUploadHandler) GetProjectFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectRepository.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Check if file exists
	filePath := filepath.Join("./static", project.File)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Get the content type
	contentType := getContentType(filePath)

	// Set appropriate headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(project.File)))
	w.Header().Set("Content-Type", contentType)

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// DownloadProjectFile forces download of the project file
func (h *FileUploadHandler) DownloadProjectFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectRepository.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Check if file exists
	filePath := filepath.Join("./static", project.File)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(project.File)))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// getContentType determines the MIME type of a file
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc", ".docx":
		return "application/msword"
	case ".txt":
		return "text/plain"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// ListUserProjects returns all projects for a specific user
func (h *FileUploadHandler) ListUserProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get all projects for this user
	projects, err := h.projectRepository.GetProjectsByUserID(userID)
	if err != nil {
		http.Error(w, "Error retrieving projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Enrich projects with metadata
	type EnrichedProject struct {
		*models.Project
		ModuleName string `json:"module_name"`
		CourseName string `json:"course_name"`
		FileType   string `json:"file_type"`
		FileSize   int64  `json:"file_size"`
	}

	enrichedProjects := make([]EnrichedProject, 0, len(projects))
	for _, project := range projects {
		// Get file information
		filePath := filepath.Join("./static", project.File)
		fileInfo, err := os.Stat(filePath)

		fileSize := int64(0)
		if err == nil {
			fileSize = fileInfo.Size()
		}

		// Get module information
		module, err := h.moduleRepository.GetModuleByID(project.ModuleID)
		moduleName := ""
		if err == nil {
			moduleName = module.Name
		}

		// You would also get course name here, but we'll use a placeholder
		courseName := "Course " + strconv.FormatInt(project.CourseID, 10)

		// Determine file type
		fileType := filepath.Ext(project.File)
		if fileType != "" {
			fileType = fileType[1:] // Remove the dot
		}

		enrichedProject := EnrichedProject{
			Project:    project,
			ModuleName: moduleName,
			CourseName: courseName,
			FileType:   fileType,
			FileSize:   fileSize,
		}

		enrichedProjects = append(enrichedProjects, enrichedProject)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrichedProjects)
}

// ListModuleProjects returns all projects for a specific module
func (h *FileUploadHandler) ListModuleProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := strconv.ParseInt(vars["moduleId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid module ID", http.StatusBadRequest)
		return
	}

	// Get all projects for this module
	projects, err := h.projectRepository.GetProjectsByModuleID(moduleID)
	if err != nil {
		http.Error(w, "Error retrieving projects: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get module information
	module, err := h.moduleRepository.GetModuleByID(moduleID)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	// Return projects with module context
	result := struct {
		Module   *models.Module    `json:"module"`
		Projects []*models.Project `json:"projects"`
	}{
		Module:   module,
		Projects: projects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeleteProject deletes a project and its associated file
func (h *FileUploadHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Get the project
	project, err := h.projectRepository.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Delete the file if it exists
	filePath := filepath.Join("./static", project.File)
	if _, err := os.Stat(filePath); err == nil {
		err = os.Remove(filePath)
		if err != nil {
			// Log error but continue with deletion from database
			fmt.Printf("Warning: Failed to delete project file: %v\n", err)
		}
	}

	// Delete project from database
	err = h.projectRepository.DeleteProject(projectID)
	if err != nil {
		http.Error(w, "Error deleting project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Also update any user modules that reference this project
	// This could be optimized by having a direct query for users with this project ID
	// but for now we'll do a targeted search using the project metadata
	userID := project.UserID
	moduleID := project.ModuleID

	if userID > 0 && moduleID > 0 {
		user, err := h.userRepository.GetUserByID(userID)
		if err == nil && user != nil {
			for i, module := range user.Modules {
				if module.ModuleID == moduleID {
					// Look for answers that reference this project
					updated := false
					for elementID, answer := range module.Answers {
						if answer.ProjectID == projectID {
							delete(module.Answers, elementID)
							updated = true
						}
					}

					if updated {
						user.Modules[i] = module
						// Update the user
						_, err = h.userRepository.UpdateUser(userID, user)
						if err != nil {
							fmt.Printf("Warning: Failed to update user after project deletion: %v\n", err)
						}
					}
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Project deleted successfully",
	})
}
