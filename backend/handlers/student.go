package controllers

import (
	"crud/config"
	"crud/models"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateStudent(w http.ResponseWriter, r *http.Request) {
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

	result, err := config.DB.Exec(query,
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

func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query("SELECT * FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		var id int
		rows.Scan(&id, &s.StudentName, &s.Address, &s.State, &s.District, &s.Taluka, &s.Gender, &s.DOB, &s.Photo, &s.Handicapped, &s.Email, &s.MobileNumber, &s.BloodGroup)
		s.ID = strconv.Itoa(id)
		students = append(students, s)
	}
	json.NewEncoder(w).Encode(students)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var s models.Student
	json.NewDecoder(r.Body).Decode(&s)

	_, err := config.DB.Exec(`UPDATE students SET studentName=?, address=?, state=?, district=?, taluka=?, gender=?, dob=?, photo=?, handicapped=?, email=?, mobileNumber=?, bloodGroup=? WHERE id=?`,
		s.StudentName, s.Address, s.State, s.District, s.Taluka, s.Gender, s.DOB, s.Photo, s.Handicapped, s.Email, s.MobileNumber, s.BloodGroup, params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Student updated")
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := config.DB.Exec("DELETE FROM students WHERE id=?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Student deleted")
}