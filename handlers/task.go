package handlers

import (
	"net/http"

	"ToDo-List/database"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	database.InsertTasks(userID, title)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idTask := r.FormValue("TaskId")
	database.DeleteTask(idTask, userID)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
