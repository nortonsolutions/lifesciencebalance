package controllers

// import User model
import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"restAPI/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// UserHandler will hold everything that controller needs
type UserHandler struct {
	userRepository models.UserRepository
	sessions       *map[string]*Session
}

// NewUserHandler returns a new UserHandler
func NewUserHandler(userRepository models.UserRepository, Sessions *map[string]*Session) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
		sessions:       Sessions,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// decode the request body into a new user struct
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.CreateIfNotExists(w, r, &user)
	h.IssueToken(w, r, &user, false)
}

func (h *UserHandler) CreateIfNotExists(w http.ResponseWriter, r *http.Request, user *models.User) {

	// See if user exists in database, if not, create
	_, err := h.userRepository.GetUserByUsername(user.Username)
	if err != nil {

		// save the user to the database
		_, err := h.userRepository.CreateUser(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Get
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// get the user from the database
	var vars = mux.Vars(r)
	var id = vars["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	user, err := h.userRepository.GetUserByID(idInt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.KeyID = idInt

	// return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

// GetAll
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// get all users from the database
	users, err := h.userRepository.GetAllUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Update
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]

	// decode the request body into a new user struct
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//convert id to int64
	idInt, _ := strconv.ParseInt(id, 10, 64)

	// update the user in the database
	_, err = h.userRepository.UpdateUser(idInt, &user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.KeyID = idInt

	// return the updated user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Delete
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// delete the user from the database
	var id = mux.Vars(r)["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	err := h.userRepository.DeleteUser(idInt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return a status notifying the client the user was deleted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (h *UserHandler) IssueToken(w http.ResponseWriter, r *http.Request, user *models.User, skipResponse bool) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(360 * time.Second)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	(*h.sessions)[sessionToken] = &Session{
		username: (*user).Username,
		expiry:   expiresAt,
	}

	type UserWithRedirect struct {
		User          models.User `json:"user"`
		Redirect_path string      `json:"redirect_path"`
	}

	userWithRedirect := UserWithRedirect{
		User:          *user,
		Redirect_path: "/app.html",
	}

	if !skipResponse {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userWithRedirect)
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.GetUserByUsernameAndPassword(creds.Username, creds.Password)

	if err != nil {
		// set status to 401
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	h.IssueToken(w, r, user, false)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var sessionToken string
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session_token" {
			sessionToken = cookie.Value
		}
	}

	if sessionToken != "" {
		delete(*h.sessions, sessionToken)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

}

func (h *UserHandler) ValidateSession(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionToken, session := GetSession(r)
		if sessionToken == "" || *session == (Session{}) || Expired(session) {
			if sessionToken != "" && *session != (Session{}) {
				delete(Sessions, sessionToken)
			}

			http.Redirect(w, r, "/expired.html", http.StatusTemporaryRedirect)
			return
		}

		expiresAt := time.Now().Add(360 * time.Second)

		session.SetTime(expiresAt)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Path:    "/",
			Value:   sessionToken,
			Expires: expiresAt,
		})

		next.ServeHTTP(w, r)

	})
}

func (h *UserHandler) GetUserByUsername(username string) (*models.User, error) {
	return h.userRepository.GetUserByUsername(username)
}

func (h *UserHandler) GetSession(r *http.Request) (string, *Session) {
	return GetSession(r)
}
