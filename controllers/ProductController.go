package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productRepository models.ProductRepository
}

func NewProductHandler(productRepository models.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepository: productRepository}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.productRepository.CreateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productRepository.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.productRepository.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.productRepository.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.productRepository.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) GetApprovedProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productRepository.GetApprovedProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetUnapprovedProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productRepository.GetUnapprovedProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	products, err := h.productRepository.GetProductsByCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductsByVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID, err := strconv.ParseInt(vars["vendorId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	products, err := h.productRepository.GetProductsByVendor(vendorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ApproveProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.productRepository.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	product.Approved = true
	_, err = h.productRepository.UpdateProduct(id, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UnapproveProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.productRepository.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	product.Approved = false
	_, err = h.productRepository.UpdateProduct(id, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
