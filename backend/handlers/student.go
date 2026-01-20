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
	(studentName, address, state, district, taluka, gender, dob, photo, handicapped, 
	email, mobileNumber, bloodGroup)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := config.DB.Exec(
		query,
		r.FormValue("studentName"),r.FormValue("address"),r.FormValue("state"),
		r.FormValue("district"),r.FormValue("taluka"),r.FormValue("gender"),
		r.FormValue("dob"),fileName,handicapped,r.FormValue("email"),
		r.FormValue("mobileNumber"),r.FormValue("bloodGroup"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := result.LastInsertId()

	session, _ := config.Store.Get(r, "student-session")

	var studentIDs []string
	if session.Values["student_ids"] != nil {
		studentIDs = session.Values["student_ids"].([]string)
	}

	studentIDs = append(studentIDs, strconv.FormatInt(lastID, 10))
	session.Values["student_ids"] = studentIDs
	session.Save(r, w)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student created and stored in session",
	})
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
		rows.Scan(
			&id,&s.StudentName,&s.Address,&s.State,&s.District,
			&s.Taluka,&s.Gender,&s.DOB,&s.Photo,&s.Handicapped,
			&s.Email,&s.MobileNumber,&s.BloodGroup,
		)
		s.ID = strconv.Itoa(id)
		students = append(students, s)
	}

	json.NewEncoder(w).Encode(students)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "student-session")

	ids, ok := session.Values["student_ids"].([]string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var s models.Student
	json.NewDecoder(r.Body).Decode(&s)

	authorized := false
	for _, id := range ids {
		if id == s.ID {
			authorized = true
			break
		}
	}

	if !authorized {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err := config.DB.Exec(`
		UPDATE students SET
		studentName=?, address=?, state=?, district=?, taluka=?, gender=?, dob=?,
		photo=?, handicapped=?, email=?, mobileNumber=?, bloodGroup=?
		WHERE id=?`,
		s.StudentName,s.Address,s.State,s.District,s.Taluka,s.Gender,
		s.DOB,s.Photo,s.Handicapped,s.Email,s.MobileNumber,s.BloodGroup,s.ID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Student updated successfully")
}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "student-session")

	ids, ok := session.Values["student_ids"].([]string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var s models.Student
	json.NewDecoder(r.Body).Decode(&s)

	authorized := false
	for _, id := range ids {
		if id == s.ID {
			authorized = true
			break
		}
	}

	if !authorized {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err := config.DB.Exec("DELETE FROM students WHERE id=?", s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updated []string
	for _, id := range ids {
		if id != s.ID {
			updated = append(updated, id)
		}
	}

	session.Values["student_ids"] = updated
	session.Save(r, w)

	json.NewEncoder(w).Encode("Student deleted successfully")
}
