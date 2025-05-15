package main

import (
	"ToDo-List/database"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var store = sessions.NewCookieStore([]byte("super-secret-key"))

type Task struct {
	ID                   int
	Title                string
	CreatedAt            time.Time
	CreatedAtAfterFormat string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := os.Getenv("connStr")
	db := database.ConnectToDb(connStr)
	defer db.Close()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/add", addTaskHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Fatal(err)
		}
		database.InsertUsers(username, password)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"]

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows := database.TaskInfo(userID)
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.CreatedAt)
		task.CreatedAtAfterFormat = task.CreatedAt.Format("02.01.2006 15:04:05")
		if err != nil {
			http.Error(w, "Ошибка при чтении задачи", 500)
			return
		}
		tasks = append(tasks, task)
	}
	slices.Reverse(tasks)

	templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Tasks": tasks,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := r.FormValue("Title")
	if title == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	database.InsertTasks(userID, title)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
