package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ModuleHandler handles product-page component operations
type ModuleHandler struct {
	moduleRepository models.ModuleRepository
	productRepository models.ProductRepository
}

// NewModuleHandler ..
func NewModuleHandler(moduleRepository models.ModuleRepository, productRepository models.ProductRepository) *ModuleHandler {
	return &ModuleHandler{moduleRepository: moduleRepository, productRepository: productRepository}
}

// add module
func (c *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	module := models.Module{}
	err := json.NewDecoder(r.Body).Decode(&module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If there was a param "productId" then use it for the product_id in module
	if id := mux.Vars(r)["productId"]; id != "" {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		module.ProductID = idInt
	}

	key, err := c.moduleRepository.CreateModule(&module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If productId is not nil, update the modules for the product with this id by adding key.ID to the modules list
	if id := mux.Vars(r)["productId"]; id != "" {
		// convert id to int64
		idInt, _ := strconv.ParseInt(id, 10, 64)
		product, err := c.productRepository.GetProductByID(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		product.Modules = append(product.Modules, key.ID)
		_, err = c.productRepository.UpdateProduct(idInt, product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key.ID)
}

// remove module
func (c *ModuleHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.moduleRepository.DeleteModule(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If productId is not nil, remove key.ID from the product modules list
	if productId := mux.Vars(r)["productId"]; productId != "" {
		// convert id to int64
		productIdInt, _ := strconv.ParseInt(productId, 10, 64)

		product, err := c.productRepository.GetProductByID(productIdInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		j := 0

		for _, module := range product.Modules {
			if module != idInt {
				product.Modules[j] = module
				j++
			}
		}

		product.Modules = product.Modules[:j]
		// product.Modules = lo.Filter(product.Modules, func(module int64, _ int) bool {
		// 	return (module != idInt)
		// })

		_, err = c.productRepository.UpdateProduct(idInt, product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Module deleted"})
}

// get all modules
func (c *ModuleHandler) GetAllModules(w http.ResponseWriter, r *http.Request) {
	modules, err := c.moduleRepository.GetAllModules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

// get module by id
func (c *ModuleHandler) GetModuleByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	module, err := c.moduleRepository.GetModuleByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(module)
}

// update module
func (c *ModuleHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	module := models.Module{}
	err := json.NewDecoder(r.Body).Decode(&module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.moduleRepository.UpdateModule(idInt, &module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

func (c *ModuleHandler) GetAllModulesByProductID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["productId"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	modules, err := c.moduleRepository.GetAllModulesByProductID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}
