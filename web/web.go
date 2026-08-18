package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed templates/* static/*
var files embed.FS

func Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/index.html", nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/login.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/register.html", nil)
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	tmpl, err := template.ParseFS(files, name)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Execute error: %v", err)
	}
}

func StaticHandler() http.Handler {
	staticFS, err := fs.Sub(files, "static")
	if err != nil {
		log.Printf("Static sub fs error: %v", err)
		return http.FileServer(http.Dir("web/static"))
	}
	return http.FileServer(http.FS(staticFS))
}