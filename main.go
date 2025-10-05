/**
 * @author Norton 2022
 * This is a simple example of a RESTful API,
 * using the gorm ORM library and the gorilla/mux router.
 *
 * Using Repository Interface per Model (best DB migration practice)
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"restAPI/routes"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// main function
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}

	jsonPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if jsonPath == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_APPLICATION_CREDENTIALS"`)
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projID, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	defer client.Close()

	// create a new router
	router := mux.NewRouter()
	routes.SetupRoutes(router, client, ctx)

	// start the server
	ssl_enabled := os.Getenv("SSL_ENABLED")
	fmt.Println("Starting the application " + projID + " @ " + time.Now().UTC().Format(time.RFC3339) + " on port " + os.Getenv("PORT") + " with SSL: " + ssl_enabled)
	if ssl_enabled == "true" {
		cert := os.Getenv("CERT_FILE")
		key := os.Getenv("KEY_FILE")
		log.Fatal(http.ListenAndServeTLS(":"+os.Getenv("PORT"), cert, key, router))
	} else {
		log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
	}
}
