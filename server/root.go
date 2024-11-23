package server

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/template/index.html")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := t.Execute(w, struct{}{}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
