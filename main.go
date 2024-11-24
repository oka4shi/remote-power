package main

import (
	"github.com/oka4shi/remote-power/server"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", server.Home)
	http.HandleFunc("POST /push", server.Push)
	http.HandleFunc("GET /push/status", server.PushStatus)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("server/template/static"))))

	log.Print(http.ListenAndServe(":8080", nil))

}
