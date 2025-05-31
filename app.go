package main

import (
	"fmt"
	"net/http"

	"proyectobd2/src/basedata"
	"proyectobd2/src/handlers"
)

func main() {
	// Paso 1: Inicializar Cassandra
	basedata.InitCassandra()

	// Paso 2: Insertar canciones si no existen
	err := basedata.SeedMusicData()
	if err != nil {
		fmt.Println("Error al insertar canciones:", err)
	}

	// Paso 3: Rutas del servidor
	http.HandleFunc("/", handlers.Index) // Página principal

	// 🚀 Rutas futuras (comentadas por ahora, pero preparadas)
	http.HandleFunc("/api/canciones", handlers.GetCanciones)
	// http.HandleFunc("/api/recomendaciones", handlers.GetRecomendaciones)
	// http.HandleFunc("/api/olap/genero", handlers.GetOLAPGenero)

	// Paso 4: Servir archivos estáticos
	fs := http.FileServer(http.Dir("src/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Paso 5: Levantar el servidor
	fmt.Println("Servidor HTTP iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
