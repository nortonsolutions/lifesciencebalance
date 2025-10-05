package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// This controller keeps track of a list of Routes
// Each Route has a name and a permission level,
// the permission level is an integer between 0 and 256,
// which is compared with the the bitmask of the user's permissions

// RouteHandler ..
type RouteHandler struct {
	RouteRepository models.RouteRepository
}

// NewRouteHandler ..
func NewRouteHandler(RouteRepository models.RouteRepository) *RouteHandler {
	return &RouteHandler{RouteRepository: RouteRepository}
}

// add Route
func (c *RouteHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	Route := models.Route{}
	err := json.NewDecoder(r.Body).Decode(&Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.RouteRepository.CreateRoute(&Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// remove Route
func (c *RouteHandler) DeleteRoute(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.RouteRepository.DeleteRoute(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Route deleted"})
}

// get all Routes
func (c *RouteHandler) GetAllRoutes(w http.ResponseWriter, r *http.Request) {
	Routes, err := c.RouteRepository.GetAllRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Routes)
}

// get Route by id
func (c *RouteHandler) GetRouteByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	Route, err := c.RouteRepository.GetRouteByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Route)
}

// update Route
func (c *RouteHandler) UpdateRoute(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	Route := models.Route{}
	err := json.NewDecoder(r.Body).Decode(&Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.RouteRepository.UpdateRoute(idInt, &Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// get Route by name
func (c *RouteHandler) GetRouteByName(w http.ResponseWriter, r *http.Request) {
	var name = mux.Vars(r)["name"]
	Route, err := c.RouteRepository.GetRouteByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Route)
}
