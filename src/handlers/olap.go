package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"

	"github.com/gin-gonic/gin"
)

func GetOLAPGenero(c *gin.Context) {
	resultados, err := basedata.GetEscuchasPorGeneroMes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos OLAP"})
		return
	}
	c.JSON(http.StatusOK, resultados)
}
