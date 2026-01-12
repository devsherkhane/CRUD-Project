package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

type Student struct {
	ID           string `json:"id"`
	StudentName  string `json:"studentName"`
	Address      string `json:"address"`
	State        string `json:"state"`
	District     string `json:"district"`
	Taluka       string `json:"taluka"`
	Gender       string `json:"gender"`
	DOB          string `json:"dob"`
	Photo        string `json:"photo"`
	Handicapped  bool   `json:"handicapped"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`
	BloodGroup   string `json:"bloodGroup"`
}

func main() {
	ConnectDB()

	// Create uploads directory if it doesn't exist
	_ = os.Mkdir("./uploads", os.ModePerm)

	router := mux.NewRouter()

	// Serve the uploads folder so the frontend can access the images
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/students", getAllStudents).Methods("GET")
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

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

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Student CRUD API</h1>"))
}
func createStudent(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("photo")
	var fileName string
	if err == nil {
		defer file.Close()
		fileName = handler.Filename
		dst, _ := os.Create(filepath.Join("./uploads", fileName))
		defer dst.Close()
		io.Copy(dst, file)
	}

	handicapped, _ := strconv.ParseBool(r.FormValue("handicapped"))

	query := `INSERT INTO students 
	(studentName, address, state, district, taluka, gender, dob, photo, handicapped, email, mobileNumber, bloodGroup)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		r.FormValue("studentName"),
		r.FormValue("address"),
		r.FormValue("state"),
		r.FormValue("district"),
		r.FormValue("taluka"),
		r.FormValue("gender"),
		r.FormValue("dob"),
		fileName,
		handicapped,
		r.FormValue("email"),
		r.FormValue("mobileNumber"),
		r.FormValue("bloodGroup"),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10), "message": "Student created"})
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := DB.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student

	for rows.Next() {
		var s Student
		var id int

		rows.Scan(
			&id,
			&s.StudentName,
			&s.Address,
			&s.State,
			&s.District,
			&s.Taluka,
			&s.Gender,
			&s.DOB,
			&s.Photo,
			&s.Handicapped,
			&s.Email,
			&s.MobileNumber,
			&s.BloodGroup,
		)

		s.ID = strconv.Itoa(id)
		students = append(students, s)
	}

	json.NewEncoder(w).Encode(students)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var s Student
	json.NewDecoder(r.Body).Decode(&s)

	_, err := DB.Exec(`
	UPDATE students SET
	studentName=?, address=?, state=?, district=?, taluka=?, gender=?, dob=?, photo=?, handicapped=?, email=?, mobileNumber=?, bloodGroup=?
	WHERE id=?`,
		s.StudentName,
		s.Address,
		s.State,
		s.District,
		s.Taluka,
		s.Gender,
		s.DOB,
		s.Photo,
		s.Handicapped,
		s.Email,
		s.MobileNumber,
		s.BloodGroup,
		params["id"],
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Student updated")
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	_, err := DB.Exec("DELETE FROM students WHERE id=?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Student deleted")
}
