package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"

	"github.com/gin-gonic/gin"
)

// GetOLAPGenero retorna las escuchas por género y mes
func GetOLAPGenero(c *gin.Context) {
	resultados, err := basedata.GetEscuchasPorGeneroMes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos OLAP"})
		return
	}
	c.JSON(http.StatusOK, resultados)
}

// GetOLAPCiudad retorna las escuchas por ciudad y mes
func GetOLAPCiudad(c *gin.Context) {
	resultados, err := basedata.GetEscuchasPorCiudadMes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos OLAP por ciudad"})
		return
	}
	c.JSON(http.StatusOK, resultados)
}
