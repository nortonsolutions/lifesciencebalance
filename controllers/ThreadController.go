package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ThreadHandler ..
type ThreadHandler struct {
	threadRepository models.ThreadRepository
	moduleRepository models.ModuleRepository
}

// NewThreadHandler ..
func NewThreadHandler(threadRepository models.ThreadRepository, moduleRepository models.ModuleRepository) *ThreadHandler {
	return &ThreadHandler{threadRepository: threadRepository, moduleRepository: moduleRepository}
}

// add thread
func (c *ThreadHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	thread := models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If there was a param "moduleId" then use it for the module_id in thread
	if id := mux.Vars(r)["moduleId"]; id != "" {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		thread.ModuleID = idInt
	}

	key, err := c.threadRepository.CreateThread(&thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If id is not nil, update the threads for the module with this id by adding key.ID to the modules list
	if id := mux.Vars(r)["moduleId"]; id != "" {
		// convert id to int64
		idInt, _ := strconv.ParseInt(id, 10, 64)
		module, err := c.moduleRepository.GetModuleByID(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		module.ThreadIDs = append(module.ThreadIDs, key.ID)
		_, err = c.moduleRepository.UpdateModule(idInt, module)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key.ID)
}

// remove thread
func (c *ThreadHandler) DeleteThread(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := c.threadRepository.DeleteThread(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If moduleId is not nil, remove key.ID from the module modules list
	if moduleId := mux.Vars(r)["moduleId"]; moduleId != "" {
		// convert id to int64
		moduleIdInt, _ := strconv.ParseInt(moduleId, 10, 64)

		module, err := c.moduleRepository.GetModuleByID(moduleIdInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		j := 0

		for _, thread := range module.ThreadIDs {
			if thread != idInt {
				module.ThreadIDs[j] = thread
				j++
			}
		}

		// module.Threads = lo.Filter(module.Threads, func(module int64, _ int) bool {
		// 	return (module != idInt)
		// })

		_, err = c.moduleRepository.UpdateModule(idInt, module)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Thread deleted"})
}

// get all threads
func (c *ThreadHandler) GetAllThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := c.threadRepository.GetAllThreads()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}

// get thread by id
func (c *ThreadHandler) GetThreadByID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	thread, err := c.threadRepository.GetThreadByID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thread)
}

// update thread
func (c *ThreadHandler) UpdateThread(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]
	thread := models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	key, err := c.threadRepository.UpdateThread(idInt, &thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// GetAllThreadsByModuleID
func (c *ThreadHandler) GetAllThreadsByModuleID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["moduleId"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	threads, err := c.threadRepository.GetAllThreadsByModuleID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threads)
}
