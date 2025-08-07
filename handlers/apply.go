package handlers

import (
	"html/template"
	"net/http"
	"studentform/database"
)

type Course struct {
	ID   int
	Name string
}

func ShowApplyForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/apply.html"))

	rows, err := database.DB.Query("SELECT id, name FROM courses")
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		rows.Scan(&c.ID, &c.Name)
		courses = append(courses, c)
	}

	data := struct {
		Courses []Course
	}{
		Courses: courses,
	}

	tmpl.Execute(w, data)
}

func ApplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		courseID := r.FormValue("course_id")
		email := r.FormValue("email")
	
		var courseName string
		err := database.DB.QueryRow("SELECT name FROM courses WHERE id = ?", courseID).Scan(&courseName)
		if err != nil {
			http.Error(w, "Failed to fetch course name: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = database.DB.Exec(`INSERT INTO applications (student_email, course_id, course_name)
		VALUES (?, ?, ?)`, email, courseID, courseName)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/apply", http.StatusSeeOther)
		return
	}
}
