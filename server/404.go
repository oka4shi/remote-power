package server

import (
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "server/template/404.html")
}
