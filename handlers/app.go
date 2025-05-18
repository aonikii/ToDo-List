package handlers

import (
	"net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["user_id"]

	if ok {
		templates.ExecuteTemplate(w, "enter.html", map[string]string{"IsAuthenticated": "Вы уже вошли в систему, можете перейти к задачам"})
		return
	}

	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "enter.html", nil)
	}
}
