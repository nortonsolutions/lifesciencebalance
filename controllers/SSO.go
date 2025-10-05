package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"restAPI/models"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssogolang *oauth2.Config
var RandomString = "asdf"

func init() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	ssogolang = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"email profile"}, // "https://www.googleapis.com/auth/userinfo.email"
		Endpoint:     google.Endpoint,
		// oauth2.Endpoint{
		// 	AuthURL: os.Getenv("AUTH_URI"), // https://accounts.google.com/o/oauth2/auth
		// 	TokenURL: os.Getenv("TOKEN_URI"), // "https://accounts.google.com/o/oauth2/token"
		// },

	}
}

// create func SSO as http middleware
func SSO(w http.ResponseWriter, r *http.Request) {
	url := ssogolang.AuthCodeURL(RandomString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *UserHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	data, err := GetUserData(state, code)
	if err != nil {
		log.Fatal(err)
	}

	// create struct to match google data
	type googleUser struct {
		ID         string `json:"id"`
		Email      string `json:"email"`
		Verified   bool   `json:"verified_email"`
		Picture    string `json:"picture"`
		Locale     string `json:"locale"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}

	gUser := googleUser{}
	err = json.Unmarshal([]byte(data), &gUser)
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username:  gUser.ID,
		Email:     gUser.Email,
		Firstname: gUser.GivenName,
		Lastname:  gUser.FamilyName,
	}

	h.CreateIfNotExists(w, r, &user)
	h.IssueToken(w, r, &user, true)

	// convert user to json string
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	// obfuscate the &user into a string
	obfuscatedUser := base64.StdEncoding.EncodeToString([]byte(userJSON))
	http.Redirect(w, r, "/app.html?user="+obfuscatedUser, http.StatusTemporaryRedirect)
}

func GetUserData(state string, code string) ([]byte, error) {
	if state != RandomString {
		return nil, fmt.Errorf("invalid state")
	}

	token, err := ssogolang.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil

}
