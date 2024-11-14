package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	"iplocate"
)

var temporalClient client.Client

// Initialize Temporal Client
func initializeTemporal() error {
	var err error
	temporalClient, err = client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	return err
}

// Start the Temporal Workflow
func startWorkflow(name string) (string, error) {
	workflowID := "workflow-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "your-task-queue",
	}

	we, err := temporalClient.ExecuteWorkflow(context.Background(), options, iplocate.GetAddressFromIP, name)
	if err != nil {
		return "", err
	}

	var result string
	err = we.Get(context.Background(), &result)
	return result, err
}

// Handle HTMX form submission
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	result, err := startWorkflow(name)
	if err != nil {
		http.Error(w, "An error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<p>%s</p>", result)
}

// Handle cURL request
func handleAPI(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := startWorkflow(requestData.Name)
	if err != nil {
		http.Error(w, "An error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := map[string]string{
		"result": result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// Serve static HTML, CSS, and JS files
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "views/index.html")
		return
	}
	http.ServeFile(w, r, filepath.Join("views", r.URL.Path))
}

func main() {
	err := initializeTemporal()
	if err != nil {
		log.Fatalf("Failed to initialize Temporal client: %v", err)
	}

	http.HandleFunc("/submit", handleSubmit)
	http.HandleFunc("/api", handleAPI)
	http.HandleFunc("/", serveStaticFiles)

	port := 4000
	fmt.Printf("Server running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
