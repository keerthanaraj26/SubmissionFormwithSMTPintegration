package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"

	_ "github.com/go-sql-driver/mysql"
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
	db  *sql.DB
)

func init() {
	var err error

	// Parse HTML template
	tpl = template.Must(template.ParseFiles("form.html"))

	// Connect to MySQL
	dsn := "root:test@123@tcp(127.0.0.1:3306)/student?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Email config
const (
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
	smtpUser     = "keerthanapushparaj25@gmail.com" // replace with your Gmail
	smtpPassword = "wxbi kbhc dcst snoq"            // use an App Password if 2FA enabled
	notifyTo     = "sivavelayutham2002@gmail.com" // email to receive notification
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

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Fill struct
	s := Student{
		FirstName: r.FormValue("first"),
		LastName:  r.FormValue("last"),
		Email:     r.FormValue("email"),
		DOB:       r.FormValue("dob"),
		Gender:    r.FormValue("gender"),
	}

	// Insert to DB
	_, err := db.Exec(
		`INSERT INTO form (firstname, lastname, email, dob, gender) VALUES (?, ?, ?, ?, ?)`,
		s.FirstName, s.LastName, s.Email, s.DOB, s.Gender)

	if err != nil {
		http.Error(w, "Database error: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	// Send email
	if err := sendEmail(s.Email, s); err != nil {
		log.Println("Email sending failed:", err)
	}

	// Redirect or confirm
	fmt.Fprintf(w, "Thanks %s! Your information has been saved and emailed.", s.FirstName)
}
