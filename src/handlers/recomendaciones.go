package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

// GetRecomendaciones retorna las canciones más escuchadas por género
func GetRecomendaciones(c *gin.Context) {
	genero := c.Query("genero")
	usuarioID := c.Query("usuario_id")
	limiteStr := c.DefaultQuery("limite", "5")
	limite, err := strconv.Atoi(limiteStr)
	if err != nil || limite <= 0 {
		limite = 5
	}

	// Si no se especifica género, buscar el favorito del usuario
	if genero == "" && usuarioID != "" {
		uid, err := gocql.ParseUUID(usuarioID)
		if err == nil {
			fav, err := basedata.GetGeneroFavoritoUsuario(uid)
			if err == nil {
				genero = fav
			}
		}
	}
	if genero == "" {
		genero = "Rock" // fallback
	}

	top, err := basedata.GetTopCancionesPorGenero(genero, limite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, top)
}
