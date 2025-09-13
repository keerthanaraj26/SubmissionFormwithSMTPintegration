package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"studentform/database"
	// "github.com/gorilla/sessions"
)

type Course struct {
	ID   int
	Name string
}
// var store = sessions.NewCookieStore([]byte("super-secret-key"))

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
		
		// session, _ := store.Get(r, "session-name")
		// session.Save(r, w)
		// emailed, ok := session.Values["email"].(string)
		// if !ok || emailed == "" {
		// 	http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
		// 	return
		// }
	
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
		if err := sendConfirmationEmail(email, courseName); err != nil {
			log.Println("Failed to send confirmation email:", err)
		}
		fmt.Fprintf(w, `Thank you! Your application for course %s has been received.`, courseName)
		
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
		return
	}
}

func sendConfirmationEmail(notifyTo, courseName string) error {
	const (
		smtpHost     = "smtp.gmail.com"
		smtpPort     = "587"
		smtpUser     = "keerthanapushparaj25@gmail.com"
		smtpPassword = "wxbi kbhc dcst snoq"            
	)

	body := fmt.Sprintf("Hello,\n\nYour application for course '%s' has been received.\n\nThank you!", courseName)

	msg := []byte("To: " + notifyTo + "\n" +
		"Subject: " + "Course Application Confirmation" + "\n\n" +
		body)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	addr := smtpHost + ":" + smtpPort

	return smtp.SendMail(addr, auth, smtpUser, []string{notifyTo}, msg)
}
