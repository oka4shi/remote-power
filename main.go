package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/oka4shi/remote-power/server"
)

var port = (func() string {
	if p := os.Getenv("PORT"); len(p) != 0 {
		return p
	} else {
		return "8000"
	}
})()

func main() {
	http.HandleFunc("GET /", server.Home)
	http.HandleFunc("POST /push", server.Push)
	http.HandleFunc("GET /push/status", server.PushStatus)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(server.TemplateDir, "/static")))))

	log.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}
