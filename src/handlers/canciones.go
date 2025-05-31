package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"

	"github.com/gin-gonic/gin"
)

func GetCanciones(c *gin.Context) {
	canciones, err := basedata.GetAllCanciones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, canciones)
}
