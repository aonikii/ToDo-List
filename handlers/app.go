package handlers

import (
	"net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "enter.html", nil)
	}
}
