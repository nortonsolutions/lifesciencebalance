package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// UserProductHandler ..
type UserProductHandler struct {
	userProductRepository models.UserProductRepository
}

// NewUserProductHandler ..
func NewUserProductHandler(userProductRepository models.UserProductRepository) *UserProductHandler {
	return &UserProductHandler{
		userProductRepository: userProductRepository,
	}
}

func (h *UserProductHandler) CreateUserProduct(w http.ResponseWriter, r *http.Request) {
	userProduct := models.UserProduct{}
	err := json.NewDecoder(r.Body).Decode(&userProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.userProductRepository.CreateUserProduct(&userProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProduct)
}

func (h *UserProductHandler) GetAllUserProducts(w http.ResponseWriter, r *http.Request) {
	userProducts, err := h.userProductRepository.GetAllUserProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProducts)
}

func (h *UserProductHandler) DeleteUserProduct(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := h.userProductRepository.DeleteUserProduct(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "UserProduct deleted"})
}

func (h *UserProductHandler) UpdateUserProduct(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	userProduct := models.UserProduct{}
	err := json.NewDecoder(r.Body).Decode(&userProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.userProductRepository.UpdateUserProduct(idInt, &userProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProduct)
}

func (h *UserProductHandler) GetUserProductByUserIDAndProductID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["userID"]
	var productID = mux.Vars(r)["productID"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	productIDInt, _ := strconv.ParseInt(productID, 10, 64)
	userProduct, err := h.userProductRepository.GetUserProductByUserIDAndProductID(userIDInt, productIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProduct)
}

func (h *UserProductHandler) GetUserProductsByProductID(w http.ResponseWriter, r *http.Request) {
	var productID = mux.Vars(r)["id"]
	productIDInt, _ := strconv.ParseInt(productID, 10, 64)
	userProducts, err := h.userProductRepository.GetUserProductsByProductID(productIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProducts)
}

func (h *UserProductHandler) GetUserProductsByUserID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	userProducts, err := h.userProductRepository.GetUserProductsByUserID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProducts)
}

func (h *UserProductHandler) GetProductsByUserID(w http.ResponseWriter, r *http.Request) {
	var userID = mux.Vars(r)["id"]
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	products, err := h.userProductRepository.GetProductsByUserID(userIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *UserProductHandler) GetUsersByProductID(w http.ResponseWriter, r *http.Request) {
	var productID = mux.Vars(r)["id"]
	productIDInt, _ := strconv.ParseInt(productID, 10, 64)
	users, err := h.userProductRepository.GetUsersByProductID(productIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserProductHandler) GetVendorsByProductID(w http.ResponseWriter, r *http.Request) {
	var productID = mux.Vars(r)["id"]
	productIDInt, _ := strconv.ParseInt(productID, 10, 64)
	users, err := h.userProductRepository.GetVendorsByProductID(productIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
