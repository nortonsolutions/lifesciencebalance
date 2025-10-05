package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderRepository    models.OrderRepository
	cartRepository     models.CartRepository
	productRepository  models.ProductRepository
	customerRepository models.CustomerRepository
}

func NewOrderHandler(orderRepository models.OrderRepository, cartRepository models.CartRepository, productRepository models.ProductRepository, customerRepository models.CustomerRepository) *OrderHandler {
	return &OrderHandler{
		orderRepository:    orderRepository,
		cartRepository:     cartRepository,
		productRepository:  productRepository,
		customerRepository: customerRepository,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = "pending"
	order.PaymentStatus = "pending"

	key, err := h.orderRepository.CreateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order.KeyID = key.ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get cart
	cart, err := h.cartRepository.GetCartByCustomerID(customerID)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	if len(cart.Items) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Get customer info
	customer, err := h.customerRepository.GetCustomerByUserID(customerID)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Convert cart items to order items and calculate total
	var orderItems []models.OrderItem
	var totalAmount float64 = 0

	for _, cartItem := range cart.Items {
		product, err := h.productRepository.GetProductByID(cartItem.ProductID)
		if err != nil {
			continue // Skip items that can't be found
		}
		orderItem := models.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(cartItem.Quantity)
	}

	// Create order
	order := models.Order{
		CustomerID:      customerID,
		Items:           orderItems,
		TotalAmount:     totalAmount,
		Status:          "pending",
		PaymentStatus:   "pending",
		ShippingAddress: customer.ShippingAddress,
		BillingAddress:  customer.BillingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	key, err := h.orderRepository.CreateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order.KeyID = key.ID

	// Clear cart after order is created
	err = h.cartRepository.ClearCart(customerID)
	if err != nil {
		// Log error but don't fail the order creation
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.orderRepository.GetOrderByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrdersByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.ParseInt(vars["customerId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orders, err := h.orderRepository.GetOrdersByCustomerID(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.orderRepository.GetOrderByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	order.Status = statusUpdate.Status
	order.UpdatedAt = time.Now()

	_, err = h.orderRepository.UpdateOrder(id, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderRepository.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
