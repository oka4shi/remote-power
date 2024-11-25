package server

import "os"

var TemplateDir = (func() string {
	dir := os.Getenv("TEMPLATE_DIR")
	if len(dir) != 0 {
		return dir
	} else {
		return "template"
	}
})()
