package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CartHandler struct {
	cartRepository models.CartRepository
}

func NewCartHandler(cartRepository models.CartRepository) *CartHandler {
	return &CartHandler{cartRepository: cartRepository}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := h.cartRepository.GetCartByCustomerID(customerID)
	if err != nil {
		// If cart doesn't exist, create a new one
		cart = &models.Cart{
			CustomerID: customerID,
			Items:      []models.CartItem{},
			UpdatedAt:  time.Now(),
		}
		_, err = h.cartRepository.CreateCart(cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item models.CartItem
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := h.cartRepository.GetCartByCustomerID(customerID)
	if err != nil {
		// Create new cart if doesn't exist
		cart = &models.Cart{
			CustomerID: customerID,
			Items:      []models.CartItem{item},
			UpdatedAt:  time.Now(),
		}
		_, err = h.cartRepository.CreateCart(cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Check if item already exists in cart
		found := false
		for i, existingItem := range cart.Items {
			if existingItem.ProductID == item.ProductID {
				cart.Items[i].Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			cart.Items = append(cart.Items, item)
		}
		cart.UpdatedAt = time.Now()
		_, err = h.cartRepository.UpdateCart(cart.KeyID, cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item models.CartItem
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := h.cartRepository.GetCartByCustomerID(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			if item.Quantity <= 0 {
				// Remove item if quantity is 0 or less
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = item.Quantity
			}
			break
		}
	}

	cart.UpdatedAt = time.Now()
	_, err = h.cartRepository.UpdateCart(cart.KeyID, cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := h.cartRepository.GetCartByCustomerID(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	cart.UpdatedAt = time.Now()
	_, err = h.cartRepository.UpdateCart(cart.KeyID, cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartRepository.ClearCart(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
