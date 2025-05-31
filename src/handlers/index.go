package handlers

import (
	"html/template"
	"net/http"
)

// Index renderiza la página principal
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/index.html")
	if err != nil {
		http.Error(w, "No se pudo cargar la plantilla: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
