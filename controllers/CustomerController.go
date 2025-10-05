package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	customerRepository models.CustomerRepository
}

func NewCustomerHandler(customerRepository models.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{customerRepository: customerRepository}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer.CreatedOn = time.Now()
	key, err := h.customerRepository.CreateCustomer(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customerRepository.GetAllCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.customerRepository.GetCustomerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetCustomerByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.customerRepository.GetCustomerByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var customer models.Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.customerRepository.UpdateCustomer(id, &customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.customerRepository.DeleteCustomer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
