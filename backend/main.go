package main

import (
	"crud/config"
	"crud/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Connect to Database
	config.ConnectDB()

	// 2. Ensure uploads directory exists
	_ = os.Mkdir("./uploads", os.ModePerm)

	// 3. Initialize Router
	router := mux.NewRouter()

	// 4. Setup Routes from the routes package
	routes.SetupRoutes(router)

	// 5. Start Server with CORS middleware
	log.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", enableCORS(router)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}