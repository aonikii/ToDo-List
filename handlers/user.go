package handlers

import (
	"ToDo-List/database"
	"net/http"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type task struct {
	ID                   int
	Title                string
	CreatedAt            time.Time
	CreatedAtAfterFormat string
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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

	var tasks []task
	for rows.Next() {
		var task task
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
