package routes

import (
	"crud/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes all the application routes
func SetupRoutes(router *mux.Router) {
	// Serve static files (photos)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Student CRUD Routes
	router.HandleFunc("/students", controllers.GetAllStudents).Methods("GET")
	router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/students/{id}", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", controllers.DeleteStudent).Methods("DELETE")
}