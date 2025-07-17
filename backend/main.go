package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	go server.WatchIfaces()

	http.Handle("GET /", http.FileServer(http.Dir(server.TemplateDir)))
	http.HandleFunc("POST /push", server.Push)
	http.HandleFunc("GET /push/status", server.PushStatus)
	http.HandleFunc("GET /network/status", server.NetworkStatus)
	http.HandleFunc("GET /network/watch", server.NetworkWatch)

	log.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}
