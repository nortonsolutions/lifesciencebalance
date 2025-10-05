package routes

import (
	"context"
	"flag"
	"net/http"
	"os"

	"restAPI/controllers"
	"restAPI/repositories"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, client *datastore.Client, ctx context.Context) {

	// Create repositories with the database connection
	userRepository := repositories.NewUserRepository(client, ctx)
	roleRepository := repositories.NewRoleRepository(client, ctx)
	routeRepository := repositories.NewRouteRepository(client, ctx)
	courseRepository := repositories.NewCourseRepository(client, ctx)
	threadRepository := repositories.NewThreadRepository(client, ctx)
	moduleRepository := repositories.NewModuleRepository(client, ctx)
	projectRepository := repositories.NewProjectRepository(client, ctx)
	userCourseRepository := repositories.NewUserCourseRepository(client, ctx)
	elementRepository := repositories.NewElementRepository(client, ctx)
	moduleElementRepository := repositories.NewModuleElementRepository(client, ctx)

	// Create handlers (controllers) with the repositories
	userHandler := controllers.NewUserHandler(userRepository, &controllers.Sessions)
	roleHandler := controllers.NewRoleHandler(roleRepository)
	routeHandler := controllers.NewRouteHandler(routeRepository)
	courseHandler := controllers.NewCourseHandler(courseRepository)
	threadHandler := controllers.NewThreadHandler(threadRepository, moduleRepository)
	moduleHandler := controllers.NewModuleHandler(moduleRepository, courseRepository)
	projectHandler := controllers.NewProjectHandler(projectRepository)
	userCourseHandler := controllers.NewUserCourseHandler(userCourseRepository)
	elementHandler := controllers.NewElementHandler(elementRepository)
	moduleElementHandler := controllers.NewModuleElementHandler(moduleElementRepository)
	progressHandler := controllers.NewProgressHandler(userRepository, courseRepository, moduleRepository, userCourseRepository)
	moduleAttemptHandler := controllers.NewModuleAttemptHandler(moduleRepository, elementRepository, moduleElementRepository, userRepository)
	fileUploadHandler := controllers.NewFileUploadHandler(projectRepository, moduleRepository, userRepository, moduleElementRepository)
	adminHandler := controllers.NewAdminHandler(userRepository, courseRepository, moduleRepository, elementRepository, projectRepository)

	// AI/Machine Learning routes
	geneticHandler := controllers.NewGeneticHandler()

	// call NewCheckPermissionsContext to create a CheckPermissionsContext
	// which will be passed to the CheckPermissions middleware
	c := NewHelperContext(routeRepository, userHandler, roleHandler)

	router.Use(c.CheckPermissions)
	// TODO: universalize the ValidateSession middleware

	// user routes - tested OK
	router.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/user", userHandler.ValidateSession(userHandler.GetAllUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.ValidateSession(userHandler.GetUser)).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.ValidateSession(userHandler.UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/{id}", userHandler.ValidateSession(userHandler.DeleteUser)).Methods("DELETE")

	// course routes - tested OK
	router.HandleFunc("/course", courseHandler.CreateCourse).Methods("POST")
	router.HandleFunc("/course", courseHandler.GetAllCourses).Methods("GET")
	router.HandleFunc("/course/{id}", courseHandler.DeleteCourse).Methods("DELETE")
	router.HandleFunc("/course/{id}", courseHandler.UpdateCourse).Methods("PUT")
	router.HandleFunc("/course/{id}", courseHandler.GetCourseByID).Methods("GET")
	router.HandleFunc("/course/{id}/instructor", userCourseHandler.GetInstructorsByCourseID).Methods("GET")

	// new course routes for approval and department filtering
	router.HandleFunc("/course/approved", courseHandler.GetApprovedCourses).Methods("GET")
	router.HandleFunc("/course/unapproved", courseHandler.GetUnapprovedCourses).Methods("GET")
	router.HandleFunc("/course/department/{department}", courseHandler.GetCoursesByDepartment).Methods("GET")
	router.HandleFunc("/course/{id}/approve", userHandler.ValidateSession(courseHandler.ApproveCourse)).Methods("PUT")
	router.HandleFunc("/course/{id}/unapprove", userHandler.ValidateSession(courseHandler.UnapproveeCourse)).Methods("PUT")

	// userCourse routes - tested OK
	router.HandleFunc("/usercourse", userCourseHandler.CreateUserCourse).Methods("POST")
	router.HandleFunc("/usercourse", userCourseHandler.GetAllUserCourses).Methods("GET")
	router.HandleFunc("/usercourse/{id}", userCourseHandler.DeleteUserCourse).Methods("DELETE")
	router.HandleFunc("/usercourse/{id}", userCourseHandler.UpdateUserCourse).Methods("PUT")
	router.HandleFunc("/usercourse/{userID}/{courseID}", userCourseHandler.GetUserCourseByUserIDAndCourseID).Methods("GET")

	// gets by UserID - tested OK
	router.HandleFunc("/user/{id}/course", userCourseHandler.GetCoursesByUserID).Methods("GET")
	router.HandleFunc("/user/{id}/usercourse", userCourseHandler.GetUserCoursesByUserID).Methods("GET")

	// gets by CourseID - tested OK
	router.HandleFunc("/course/{id}/user", userCourseHandler.GetUsersByCourseID).Methods("GET")
	router.HandleFunc("/course/{id}/usercourse", userCourseHandler.GetUserCoursesByCourseID).Methods("GET")

	// role routes - tested OK
	router.HandleFunc("/role", roleHandler.CreateRole).Methods("POST")
	router.HandleFunc("/role", roleHandler.GetAllRoles).Methods("GET")
	router.HandleFunc("/role/{id}", roleHandler.DeleteRole).Methods("DELETE")
	router.HandleFunc("/role/{id}", roleHandler.UpdateRole).Methods("PUT")
	router.HandleFunc("/role/{id}", roleHandler.GetRoleByID).Methods("GET")

	// route routes - tested OK
	router.HandleFunc("/route", routeHandler.CreateRoute).Methods("POST")
	router.HandleFunc("/route", routeHandler.GetAllRoutes).Methods("GET")
	router.HandleFunc("/route/{id}", routeHandler.DeleteRoute).Methods("DELETE")
	router.HandleFunc("/route/{id}", routeHandler.UpdateRoute).Methods("PUT")

	// module routes - tested OK
	router.HandleFunc("/module", moduleHandler.GetAllModules).Methods("GET")
	router.HandleFunc("/module/{id}", moduleHandler.UpdateModule).Methods("PUT")
	router.HandleFunc("/module/{id}", moduleHandler.GetModuleByID).Methods("GET")

	// course-module routes - tested OK
	router.HandleFunc("/course/{courseId}/module", moduleHandler.GetAllModulesByCourseID).Methods("GET")
	router.HandleFunc("/course/{courseId}/module", moduleHandler.CreateModule).Methods("POST")
	router.HandleFunc("/course/{courseId}/module/{id}", moduleHandler.DeleteModule).Methods("DELETE")
	router.HandleFunc("/course/{courseId}/module/{id}", moduleHandler.UpdateModule).Methods("PUT")
	router.HandleFunc("/course/{courseId}/module/{id}", moduleHandler.GetModuleByID).Methods("GET")

	// thread routes - tested OK
	router.HandleFunc("/thread", threadHandler.GetAllThreads).Methods("GET")
	router.HandleFunc("/thread/{id}", threadHandler.UpdateThread).Methods("PUT")
	router.HandleFunc("/thread/{id}", threadHandler.GetThreadByID).Methods("GET")

	// module-thread routes - tested OK
	router.HandleFunc("/module/{moduleId}/thread", threadHandler.CreateThread).Methods("POST")
	router.HandleFunc("/module/{moduleId}/thread", threadHandler.GetAllThreadsByModuleID).Methods("GET")
	router.HandleFunc("/module/{moduleId}/thread/{id}", threadHandler.DeleteThread).Methods("DELETE")
	router.HandleFunc("/module/{moduleId}/thread/{id}", threadHandler.UpdateThread).Methods("PUT")
	router.HandleFunc("/module/{moduleId}/thread/{id}", threadHandler.GetThreadByID).Methods("GET")

	// element routes - tested OK
	router.HandleFunc("/element", elementHandler.CreateElement).Methods("POST")
	router.HandleFunc("/element", elementHandler.GetAllElements).Methods("GET")
	router.HandleFunc("/element/{id}", elementHandler.DeleteElement).Methods("DELETE")
	router.HandleFunc("/element/{id}", elementHandler.UpdateElement).Methods("PUT")
	router.HandleFunc("/element/{id}", elementHandler.GetElementByID).Methods("GET")

	// moduleElement routes - tested OK
	router.HandleFunc("/moduleelement", moduleElementHandler.CreateModuleElement).Methods("POST")
	router.HandleFunc("/moduleelement", moduleElementHandler.GetAllModuleElements).Methods("GET")
	router.HandleFunc("/moduleelement/{id}", moduleElementHandler.DeleteModuleElement).Methods("DELETE")
	router.HandleFunc("/moduleelement/{id}", moduleElementHandler.UpdateModuleElement).Methods("PUT")
	router.HandleFunc("/moduleelement/{moduleID}/{elementID}", moduleElementHandler.GetModuleElementByModuleIDAndElementID).Methods("GET")

	// gets by ModuleID - tested OK
	router.HandleFunc("/module/{id}/element", moduleElementHandler.GetElementsByModuleID).Methods("GET")
	router.HandleFunc("/module/{id}/moduleelement", moduleElementHandler.GetModuleElementsByModuleID).Methods("GET")

	// gets by ElementID - tested OK
	router.HandleFunc("/element/{id}/module", moduleElementHandler.GetModulesByElementID).Methods("GET")
	router.HandleFunc("/element/{id}/moduleelement", moduleElementHandler.GetModuleElementsByElementID).Methods("GET")

	// project routes - TODO: test
	router.HandleFunc("/project", projectHandler.CreateProject).Methods("POST")
	router.HandleFunc("/project", projectHandler.GetAllProjects).Methods("GET")
	router.HandleFunc("/project/{id}", projectHandler.DeleteProject).Methods("DELETE")
	router.HandleFunc("/project/{id}", projectHandler.UpdateProject).Methods("PUT")
	router.HandleFunc("/project/{id}", projectHandler.GetProjectByID).Methods("GET")
	router.HandleFunc("/user/{id}/project", projectHandler.GetProjectsByUserID).Methods("GET")

	// login/auth routes
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/sso", controllers.SSO).Methods("GET")
	router.HandleFunc("/callback", userHandler.Callback).Methods("GET")
	router.HandleFunc("/logout", userHandler.Logout).Methods("GET")

	// genetic algorithm routes - TODO: test
	router.HandleFunc("/genetic", geneticHandler.RunGenetic).Methods("POST")

	// admin routes
	router.HandleFunc("/admin/stats", userHandler.ValidateSession(adminHandler.GetSystemStats)).Methods("GET")
	router.HandleFunc("/admin/user/{id}/stats", userHandler.ValidateSession(adminHandler.GetUserStats)).Methods("GET")
	router.HandleFunc("/admin/course/{id}/stats", userHandler.ValidateSession(adminHandler.GetCourseStats)).Methods("GET")

	// progress tracking routes
	router.HandleFunc("/user/{userId}/progress", userHandler.ValidateSession(progressHandler.GetUserProgress)).Methods("GET")
	router.HandleFunc("/user/{userId}/course/{courseId}/progress", userHandler.ValidateSession(progressHandler.GetCourseProgress)).Methods("GET")
	router.HandleFunc("/user/{userId}/course/{courseId}/progress", userHandler.ValidateSession(progressHandler.UpdateUserCourseProgress)).Methods("PUT")

	// module attempt routes
	router.HandleFunc("/module/{id}/start", userHandler.ValidateSession(moduleAttemptHandler.StartModule)).Methods("GET")
	router.HandleFunc("/user/{userId}/module/{id}/submit", userHandler.ValidateSession(moduleAttemptHandler.SubmitModule)).Methods("POST")
	router.HandleFunc("/user/{userId}/module/{id}/results", userHandler.ValidateSession(moduleAttemptHandler.GetModuleResults)).Methods("GET")
	router.HandleFunc("/module/{id}/analytics", userHandler.ValidateSession(moduleAttemptHandler.GetModuleAnalytics)).Methods("GET")
	router.HandleFunc("/user/{userId}/module/{id}/reset", userHandler.ValidateSession(moduleAttemptHandler.ResetModuleAttempt)).Methods("POST")

	// file upload routes
	router.HandleFunc("/upload/project", userHandler.ValidateSession(fileUploadHandler.UploadProject)).Methods("POST")
	router.HandleFunc("/project/{id}/file", userHandler.ValidateSession(fileUploadHandler.GetProjectFile)).Methods("GET")
	router.HandleFunc("/project/{id}/download", userHandler.ValidateSession(fileUploadHandler.DownloadProjectFile)).Methods("GET")
	router.HandleFunc("/project/{id}", userHandler.ValidateSession(fileUploadHandler.DeleteProject)).Methods("DELETE")
	router.HandleFunc("/user/{userId}/projects", userHandler.ValidateSession(fileUploadHandler.ListUserProjects)).Methods("GET")
	router.HandleFunc("/module/{moduleId}/projects", userHandler.ValidateSession(fileUploadHandler.ListModuleProjects)).Methods("GET")

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}

	var dir string // command-line flag will override environment variable
	flag.StringVar(&dir, "dir", staticDir, "the directory to serve files from")
	flag.Parse()

	// This will serve files under http://localhost:8000/<filename> in the 'dir' directory.
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	// TODO: cache these routes!
	c.UpdateRoutes(router)

}
