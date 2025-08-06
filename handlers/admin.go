package handlers

import (
	"html/template"
	"net/http"
	"net/smtp"
	"studentform/database"
	"fmt"
	"log"
)

type Student struct {
	FirstName string
	LastName  string
	Email     string
	DOB       string
	Gender    string
}

var (

	tpl *template.Template
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseFiles("form.html"))

	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUser     = "keerthanapushparaj25@gmail.com"
	smtpPassword = "wxbi kbhc dcst snoq"            
)

func sendEmail(notifyTo string, s Student) error {
	body := fmt.Sprintf(
		`Thank you for registering the student form.

Here are your submitted details:
First Name: %s
Last Name:  %s
Email:      %s
DOB:        %s
Gender:     %s,

Regards,
Your Team`,
		s.FirstName, s.LastName, s.Email, s.DOB, s.Gender)

	msg := []byte("To: " + notifyTo + "\r\n" +
		"Subject: New Student Registration\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	addr := smtpHost + ":" + smtpPort

	return smtp.SendMail(addr, auth, smtpUser, []string{notifyTo}, msg)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	s := Student{
		FirstName: r.FormValue("first"),
		LastName:  r.FormValue("last"),
		Email:     r.FormValue("email"),
		DOB:       r.FormValue("dob"),
		Gender:    r.FormValue("gender"),
	}

	_, err := database.DB.Exec(
		`INSERT INTO form (firstname, lastname, email, dob, gender) VALUES (?, ?, ?, ?, ?)`,
		s.FirstName, s.LastName, s.Email, s.DOB, s.Gender)

	if err != nil {
		http.Error(w, "Database error: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendEmail(s.Email, s); err != nil {
		log.Println("Email sending failed:", err)
	}

	fmt.Fprintf(w, "Thanks %s! Your information has been saved and emailed.", s.FirstName)
}
