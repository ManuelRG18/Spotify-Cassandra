package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type EscuchaInput struct {
	UsuarioID string `json:"usuario_id"`
	CancionID string `json:"cancion_id"`
	Fecha     string `json:"fecha_escucha"` // formato YYYY-MM-DD
}

func RegistrarEscucha(c *gin.Context) {
	var input EscuchaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	usuarioID, err := gocql.ParseUUID(input.UsuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de usuario inválido"})
		return
	}
	cancionID, err := gocql.ParseUUID(input.CancionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de canción inválido"})
		return
	}

	err = basedata.RegistrarEscucha(usuarioID, cancionID, input.Fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Escucha registrada correctamente"})
}
