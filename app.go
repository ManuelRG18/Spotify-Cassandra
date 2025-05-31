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

		router := gin.Default()

		router.Static("/static", "./src/static")

		router.LoadHTMLGlob("src/templates/*")

		router.GET("/", handlers.Index)

		api := router.Group("/api")
		{
			api.GET("/canciones", handlers.GetCanciones)
			// api.POST("/usuarios", handlers.CreateUsuario)
			// api.POST("/escuchar", handlers.RegistrarEscucha)
			// api.GET("/recomendaciones", handlers.GetRecomendaciones)
			// api.GET("/olap/genero", handlers.GetOLAPGenero)
		}

		// Iniciar servidor
		router.Run(":8080")
	}

}
