package main

import (
	"fmt"

	"proyectobd2/src/basedata"
	"proyectobd2/src/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Paso 1: Inicializar Cassandra
	basedata.InitCassandra()

	// Paso 2: Insertar canciones si no existen
	err := basedata.SeedMusicData()
	if err != nil {
		fmt.Println("Error al insertar canciones:", err)
	}

	// 🔄 Mover esto fuera del if
	router := gin.Default()

	// Servir archivos estáticos
	router.Static("/static", "./src/static")

	// Cargar plantillas HTML
	router.LoadHTMLGlob("src/templates/*")

	// Ruta principal
	router.GET("/", handlers.Index)

	// Rutas de API
	api := router.Group("/api")
	{
		api.GET("/canciones", handlers.GetCanciones)
		// api.POST("/usuarios", handlers.CreateUsuario)
		// api.POST("/escuchar", handlers.RegistrarEscucha)
		// api.GET("/recomendaciones", handlers.GetRecomendaciones)
		// api.GET("/olap/genero", handlers.GetOLAPGenero)
	}

	// Iniciar servidor
	fmt.Println("Servidor iniciado en http://localhost:8080")
	router.Run(":8080")
}
