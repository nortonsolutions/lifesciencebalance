package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type VendorHandler struct {
	vendorRepository models.VendorRepository
}

func NewVendorHandler(vendorRepository models.VendorRepository) *VendorHandler {
	return &VendorHandler{vendorRepository: vendorRepository}
}

func (h *VendorHandler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var vendor models.Vendor
	err := json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vendor.CreatedOn = time.Now()
	vendor.Approved = false // Default to unapproved
	key, err := h.vendorRepository.CreateVendor(&vendor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vendor.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) GetAllVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.vendorRepository.GetAllVendors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}

func (h *VendorHandler) GetVendorByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vendor, err := h.vendorRepository.GetVendorByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) GetVendorByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vendor, err := h.vendorRepository.GetVendorByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var vendor models.Vendor
	err = json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.vendorRepository.UpdateVendor(id, &vendor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vendor.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

func (h *VendorHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.vendorRepository.DeleteVendor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VendorHandler) GetApprovedVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.vendorRepository.GetApprovedVendors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}

func (h *VendorHandler) ApproveVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vendor, err := h.vendorRepository.GetVendorByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vendor.Approved = true
	_, err = h.vendorRepository.UpdateVendor(id, vendor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}
