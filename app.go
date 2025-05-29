package main

import (
	"fmt"
	"html/template"
	"net/http"
	basedata "proyectobd2/src/Basedata"
)

func index(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/index.html")
	if err != nil {
		http.Error(rw, "No se pudo cargar la plantilla: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(rw, nil)
}

func main() {
	// Inicializar Cassandra y crear tablas
	basedata.InitCassandra()

	// Insertar canciones de ejemplo si no existen
	err := basedata.SeedMusicData()
	if err != nil {
		fmt.Println("Error al insertar canciones:", err)
	}

	http.HandleFunc("/", index)

	fmt.Println("Servidor HTTP iniciado en http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
