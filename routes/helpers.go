package routes

import (
	"net/http"
	"restAPI/controllers"
	"restAPI/models"
	"restAPI/repositories"

	"github.com/gorilla/mux"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// "Receiver" model which allows main function to pass
// routeRepository, userHandler, and roleHandler to the context
type HelperContext struct {
	routeRepository *repositories.BaseRepository
	userHandler     *controllers.UserHandler
	roleHandler     *controllers.RoleHandler
}

// Function to return a new context
func NewHelperContext(routeRepository *repositories.BaseRepository, userHandler *controllers.UserHandler, roleHandler *controllers.RoleHandler) *HelperContext {
	return &HelperContext{
		routeRepository: routeRepository,
		userHandler:     userHandler,
		roleHandler:     roleHandler,
	}
}

func (ctx *HelperContext) UpdateRoutes(router *mux.Router) {
	// Add the name of each route to an array
	var routeNames []string
	routesInRepository, _ := ctx.routeRepository.GetAllRoutes()
	for _, r := range routesInRepository {
		routeNames = append(routeNames, r.Name)
	}

	var routesConfigured []string
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		name, _ := route.GetPathTemplate()
		method, err2 := route.GetMethods()
		if err2 == nil {
			name += "_" + method[0]
			// fmt.Println(name)
			routesConfigured = append(routesConfigured, name)
		}
		return nil
	})

	// Check if any routesConfigured are missing from routeNames
	// If so, add them to the database
	for _, r := range routesConfigured {
		if !Contains(routeNames, r) {
			ctx.routeRepository.CreateRoute(&models.Route{Name: r, PermissionLevel: 0})
		}
	}
}

// Middleware to check user roles
func (ctx *HelperContext) CheckPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Does the resource require authorization?
		route := mux.CurrentRoute(r)
		name, _ := route.GetPathTemplate()
		method, err2 := route.GetMethods()
		if err2 == nil {
			name += "_" + method[0]

			// Get the route from the database
			router, _ := ctx.routeRepository.GetRouteByName(name)
			if router != nil && router.PermissionLevel > 0 {
				// Get the user from the session
				_, session := ctx.userHandler.GetSession(r)

				if *session == (controllers.Session{}) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				user, _ := ctx.userHandler.GetUserByUsername(session.GetUsername())
				roleKey := ctx.roleHandler.GetRoleKey(user.GetRoles())

				if user == nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// Compare bitwise AND of the roleKey and the required permission level
				if (roleKey & router.PermissionLevel) == 0 {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

			}
		}

		next.ServeHTTP(w, r)

		// do something after the request
	})
}
