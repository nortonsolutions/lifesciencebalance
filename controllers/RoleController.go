package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// This controller keeps track of a list of roles
// Each role is mapped to a specific bit in a bitmask
// Each user will have a corresponding bitmask to determine which roles they have
// This allows us to easily check if a user has a specific role

// RoleHandler ..
type RoleHandler struct {
	roleRepository models.RoleRepository
}

// NewRoleHandler ..
func NewRoleHandler(roleRepository models.RoleRepository) *RoleHandler {
	return &RoleHandler{roleRepository: roleRepository}
}

// add role
func (c *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	role := models.Role{}
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.roleRepository.CreateRole(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// remove role
func (c *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.roleRepository.DeleteRole(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Role deleted"})
}

// get all roles
func (c *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := c.roleRepository.GetAllRoles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// get role by id
func (c *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	role, err := c.roleRepository.GetRoleByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// update role
func (c *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	role := models.Role{}
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.roleRepository.UpdateRole(idInt, &role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// get role by name
func (c *RoleHandler) GetRoleByName(w http.ResponseWriter, r *http.Request) {
	var name = mux.Vars(r)["name"]
	role, err := c.roleRepository.GetRoleByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (c *RoleHandler) GetRoleKey(roles []string) int {

	var roleKey int
	for _, userRole := range roles {
		// What is the NumericValue for this role?
		role, _ := c.roleRepository.GetRoleByName(userRole)
		roleKey += role.NumericValue
	}
	return roleKey

}
