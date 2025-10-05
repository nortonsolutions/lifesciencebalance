package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ModuleElementHandler ..
type ModuleElementHandler struct {
	moduleElementRepository models.ModuleElementRepository
}

// NewModuleElementHandler ..
func NewModuleElementHandler(moduleElementRepository models.ModuleElementRepository) *ModuleElementHandler {
	return &ModuleElementHandler{
		moduleElementRepository: moduleElementRepository,
	}
}

func (h *ModuleElementHandler) CreateModuleElement(w http.ResponseWriter, r *http.Request) {
	moduleElement := models.ModuleElement{}
	err := json.NewDecoder(r.Body).Decode(&moduleElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.moduleElementRepository.CreateModuleElement(&moduleElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	moduleElement.KeyID = key.ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElement)
}

func (h *ModuleElementHandler) GetAllModuleElements(w http.ResponseWriter, r *http.Request) {
	moduleElements, err := h.moduleElementRepository.GetAllModuleElements()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElements)
}

func (h *ModuleElementHandler) DeleteModuleElement(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := h.moduleElementRepository.DeleteModuleElement(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "ModuleElement deleted"})
}

func (h *ModuleElementHandler) UpdateModuleElement(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	moduleElement := models.ModuleElement{}
	err := json.NewDecoder(r.Body).Decode(&moduleElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.moduleElementRepository.UpdateModuleElement(idInt, &moduleElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	moduleElement.KeyID = key.ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElement)
}

func (h *ModuleElementHandler) GetModuleElementByModuleIDAndElementID(w http.ResponseWriter, r *http.Request) {
	var moduleID = mux.Vars(r)["moduleID"]
	var elementID = mux.Vars(r)["elementID"]
	moduleIDInt, _ := strconv.ParseInt(moduleID, 10, 64)
	elementIDInt, _ := strconv.ParseInt(elementID, 10, 64)
	moduleElement, err := h.moduleElementRepository.GetModuleElementByModuleIDAndElementID(moduleIDInt, elementIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElement)
}

func (h *ModuleElementHandler) GetModuleElementsByElementID(w http.ResponseWriter, r *http.Request) {
	var elementID = mux.Vars(r)["id"]
	elementIDInt, _ := strconv.ParseInt(elementID, 10, 64)
	moduleElements, err := h.moduleElementRepository.GetModuleElementsByElementID(elementIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElements)
}

func (h *ModuleElementHandler) GetModuleElementsByModuleID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	moduleElements, err := h.moduleElementRepository.GetModuleElementsByModuleID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moduleElements)
}

func (h *ModuleElementHandler) GetElementsByModuleID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	elements, err := h.moduleElementRepository.GetElementsByModuleID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(elements)
}

func (h *ModuleElementHandler) GetModulesByElementID(w http.ResponseWriter, r *http.Request) {
	var elementID = mux.Vars(r)["id"]
	elementIDInt, _ := strconv.ParseInt(elementID, 10, 64)
	modules, err := h.moduleElementRepository.GetModulesByElementID(elementIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

func (h *ModuleElementHandler) GetElementsByInstructorID(w http.ResponseWriter, r *http.Request) {
	var elementID = mux.Vars(r)["id"]
	elementIDInt, _ := strconv.ParseInt(elementID, 10, 64)
	users, err := h.moduleElementRepository.GetElementsByInstructorID(elementIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
