package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRecomendaciones retorna las canciones más escuchadas por género
func GetRecomendaciones(c *gin.Context) {
	genero := c.DefaultQuery("genero", "Rock") // Por defecto Rock
	limiteStr := c.DefaultQuery("limite", "5")
	limite, err := strconv.Atoi(limiteStr)
	if err != nil || limite <= 0 {
		limite = 5
	}

	top, err := basedata.GetTopCancionesPorGenero(genero, limite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, top)
}
