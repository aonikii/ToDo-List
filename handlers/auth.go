package handlers

import (
	"ToDo-List/database"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "register.html", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			templates.ExecuteTemplate(w, "register.html", map[string]string{"Error": "Заполните все поля"})
			return
		}
		password, err := hashPassword(password)
		if err != nil {
			log.Panic(err)
		}
		database.InsertUsers(username, password)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "login.html", nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		var userID int
		userID, dbPassword, err := database.LoginCheck(username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Неверное имя пользователя или пароль"})
			return
		}

		// Сохраняем userID в сессию
		session, _ := store.Get(r, "session")
		session.Values["user_id"] = userID
		session.Save(r, w)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
