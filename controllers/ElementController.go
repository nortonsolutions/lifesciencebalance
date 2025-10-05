package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ElementHandler ..
type ElementHandler struct {
	elementRepository models.ElementRepository
}

// NewElementHandler ..
func NewElementHandler(elementRepository models.ElementRepository) *ElementHandler {
	return &ElementHandler{elementRepository: elementRepository}
}

// add element
func (c *ElementHandler) CreateElement(w http.ResponseWriter, r *http.Request) {
	element := models.Element{}
	err := json.NewDecoder(r.Body).Decode(&element)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.elementRepository.CreateElement(&element)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key.ID)
}

// remove element
func (c *ElementHandler) DeleteElement(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.elementRepository.DeleteElement(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Element deleted"})
}

// get all elements
func (c *ElementHandler) GetAllElements(w http.ResponseWriter, r *http.Request) {
	elements, err := c.elementRepository.GetAllElements()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(elements)
}

// get element by id
func (c *ElementHandler) GetElementByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	element, err := c.elementRepository.GetElementByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(element)
}

// update element
func (c *ElementHandler) UpdateElement(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	element := models.Element{}
	err := json.NewDecoder(r.Body).Decode(&element)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.elementRepository.UpdateElement(idInt, &element)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}
