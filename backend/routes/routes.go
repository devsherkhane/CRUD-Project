package routes

import (
	"crud/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.HandleFunc("/",home).Methods("GET")
	router.HandleFunc("/students", controllers.GetAllStudents).Methods("GET")
	router.HandleFunc("/students", controllers.CreateStudent).Methods("POST")
	router.HandleFunc("/students", controllers.UpdateStudent).Methods("PUT")
	router.HandleFunc("/students", controllers.DeleteStudent).Methods("DELETE")

	router.HandleFunc("/students/pdf", controllers.DownloadStudentsPDF).Methods("GET")

}
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome CRUD Project by Dev</h1>"))
}