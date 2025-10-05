package controllers

import (
	"encoding/json"
	"net/http"
	"restAPI/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ModuleHandler ..
type ModuleHandler struct {
	moduleRepository models.ModuleRepository
	courseRepository models.CourseRepository
}

// NewModuleHandler ..
func NewModuleHandler(moduleRepository models.ModuleRepository, courseRepository models.CourseRepository) *ModuleHandler {
	return &ModuleHandler{moduleRepository: moduleRepository, courseRepository: courseRepository}
}

// add module
func (c *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	module := models.Module{}
	err := json.NewDecoder(r.Body).Decode(&module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If there was a param "courseId" then use it for the course_id in module
	if id := mux.Vars(r)["courseId"]; id != "" {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		module.CourseID = idInt
	}

	key, err := c.moduleRepository.CreateModule(&module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If courseId is not nil, update the modules for the course with this id by adding key.ID to the modules list
	if id := mux.Vars(r)["courseId"]; id != "" {
		// convert id to int64
		idInt, _ := strconv.ParseInt(id, 10, 64)
		course, err := c.courseRepository.GetCourseByID(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		course.Modules = append(course.Modules, key.ID)
		_, err = c.courseRepository.UpdateCourse(idInt, course)
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

	// If courseId is not nil, remove key.ID from the course modules list
	if courseId := mux.Vars(r)["courseId"]; courseId != "" {
		// convert id to int64
		courseIdInt, _ := strconv.ParseInt(courseId, 10, 64)

		course, err := c.courseRepository.GetCourseByID(courseIdInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		j := 0

		for _, module := range course.Modules {
			if module != idInt {
				course.Modules[j] = module
				j++
			}
		}

		course.Modules = course.Modules[:j]
		// course.Modules = lo.Filter(course.Modules, func(module int64, _ int) bool {
		// 	return (module != idInt)
		// })

		_, err = c.courseRepository.UpdateCourse(idInt, course)
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

func (c *ModuleHandler) GetAllModulesByCourseID(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["courseId"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	modules, err := c.moduleRepository.GetAllModulesByCourseID(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}
