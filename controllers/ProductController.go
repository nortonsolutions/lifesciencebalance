package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ProductHandler ..
type ProductHandler struct {
	productRepository models.ProductRepository
}

// NewProductHandler ..
func NewProductHandler(productRepository models.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepository: productRepository}
}

// add product
func (c *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.productRepository.CreateProduct(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key.ID)
}

// remove product
func (c *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.productRepository.DeleteProduct(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
}

// get all products
func (c *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productRepository.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// get product by id
func (c *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	product, err := c.productRepository.GetProductByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// update product
func (c *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	product := models.Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.productRepository.UpdateProduct(idInt, &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// GetApprovedProducts returns only approved products
func (c *ProductHandler) GetApprovedProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productRepository.GetApprovedProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetUnapprovedProducts returns only unapproved products
func (c *ProductHandler) GetUnapprovedProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productRepository.GetUnapprovedProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProductsByDepartment returns products filtered by department
func (c *ProductHandler) GetProductsByDepartment(w http.ResponseWriter, r *http.Request) {
	var department = mux.Vars(r)["department"]
	products, err := c.productRepository.GetProductsByDepartment(department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// ApproveProduct approves a product
func (c *ProductHandler) ApproveProduct(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	// Get the product first
	product, err := c.productRepository.GetProductByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set approved to true
	product.Approved = true

	// Update the product
	key, err := c.productRepository.UpdateProduct(idInt, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// UnapproveProduct unapproves a product
func (c *ProductHandler) UnapproveProduct(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)

	// Get the product first
	product, err := c.productRepository.GetProductByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set approved to false
	product.Approved = false

	// Update the product
	key, err := c.productRepository.UpdateProduct(idInt, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}
