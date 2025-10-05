package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// UserCourseHandler ..
type UserCourseHandler struct {
	userCourseRepository models.UserCourseRepository
}

// NewUserCourseHandler ..
func NewUserCourseHandler(userCourseRepository models.UserCourseRepository) *UserCourseHandler {
	return &UserCourseHandler{
		userCourseRepository: userCourseRepository,
	}
}

func (h *UserCourseHandler) CreateUserCourse(w http.ResponseWriter, r *http.Request) {
	userCourse := models.UserCourse{}
	err := json.NewDecoder(r.Body).Decode(&userCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.userCourseRepository.CreateUserCourse(&userCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourse)
}

func (h *UserCourseHandler) GetAllUserCourses(w http.ResponseWriter, r *http.Request) {
	userCourses, err := h.userCourseRepository.GetAllUserCourses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourses)
}

func (h *UserCourseHandler) DeleteUserCourse(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := h.userCourseRepository.DeleteUserCourse(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "UserCourse deleted"})
}

func (h *UserCourseHandler) UpdateUserCourse(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	userCourse := models.UserCourse{}
	err := json.NewDecoder(r.Body).Decode(&userCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.userCourseRepository.UpdateUserCourse(idInt, &userCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourse)
}

func (h *UserCourseHandler) GetUserCourseByUserIDAndCourseID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["userID"]
	var courseID = mux.Vars(r)["courseID"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	courseIDInt, _ := strconv.ParseInt(courseID, 10, 64)
	userCourse, err := h.userCourseRepository.GetUserCourseByUserIDAndCourseID(userIDInt, courseIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourse)
}

func (h *UserCourseHandler) GetUserCoursesByCourseID(w http.ResponseWriter, r *http.Request) {
	var courseID = mux.Vars(r)["id"]
	courseIDInt, _ := strconv.ParseInt(courseID, 10, 64)
	userCourses, err := h.userCourseRepository.GetUserCoursesByCourseID(courseIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourses)
}

func (h *UserCourseHandler) GetUserCoursesByUserID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	userCourses, err := h.userCourseRepository.GetUserCoursesByUserID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCourses)
}

func (h *UserCourseHandler) GetCoursesByUserID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	courses, err := h.userCourseRepository.GetCoursesByUserID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *UserCourseHandler) GetUsersByCourseID(w http.ResponseWriter, r *http.Request) {
	var courseID = mux.Vars(r)["id"]
	courseIDInt, _ := strconv.ParseInt(courseID, 10, 64)
	users, err := h.userCourseRepository.GetUsersByCourseID(courseIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserCourseHandler) GetInstructorsByCourseID(w http.ResponseWriter, r *http.Request) {
	var courseID = mux.Vars(r)["id"]
	courseIDInt, _ := strconv.ParseInt(courseID, 10, 64)
	users, err := h.userCourseRepository.GetInstructorsByCourseID(courseIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
