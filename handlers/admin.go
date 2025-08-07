package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"studentform/database"
	

	"golang.org/x/crypto/bcrypt"
)

type Student struct {
	FirstName string
	LastName  string
	Email     string
	DOB       string
	Gender    string
	Password  string
}

var (
	tpl *template.Template
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseFiles("templates/register.html"))
	tpl.Execute(w, nil)
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
	Your Team`, s.FirstName, s.LastName, s.Email, s.DOB, s.Gender)

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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		Password:  r.FormValue("pswd"),
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec(
		`INSERT INTO form (firstname, lastname, email, dob, gender, password_hash) VALUES (?, ?, ?, ?, ?, ?)`,
		s.FirstName, s.LastName, s.Email, s.DOB, s.Gender, hashedPassword)

	if err != nil {
		http.Error(w, "Database error: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendEmail(s.Email, s); err != nil {
		log.Println("Email sending failed:", err)
	}
	log.Printf(`Thanks %s! Your information has been saved and emailed.`, s.FirstName)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var pass string
		err := database.DB.QueryRow("SELECT password_hash FROM form WHERE email = ?", email).Scan(&pass)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}