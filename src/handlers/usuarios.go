package handlers

import (
	"net/http"
	"proyectobd2/src/basedata"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type UsuarioInput struct {
	Nombre   string `json:"nombre"`
	Ciudad   string `json:"ciudad"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUsuario(c *gin.Context) {
	var input UsuarioInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	id, err := basedata.InsertUsuario(input.Nombre, input.Ciudad, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje":    "Usuario registrado correctamente",
		"usuario_id": id,
	})
}

func LoginUsuario(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	var id gocql.UUID
	var nombre, passwordDB string
	err := basedata.Session.Query(`SELECT id, nombre, password FROM usuarios WHERE email = ? LIMIT 1`,
		input.Email).Scan(&id, &nombre, &passwordDB)

	if err != nil || passwordDB != input.Password {
		c.JSON(401, gin.H{"error": "Correo o contraseña incorrectos"})
		return
	}

	c.JSON(200, gin.H{"usuario_id": id, "nombre": nombre})
}

// Obtener datos de un usuario por ID
func GetUsuarioByID(c *gin.Context) {
	idStr := c.Param("id")
	var nombre, ciudad string
	var id gocql.UUID
	var email string
	err := basedata.Session.Query(`SELECT id, nombre, ciudad, email FROM usuarios WHERE id = ? LIMIT 1`, idStr).Scan(&id, &nombre, &ciudad, &email)
	if err != nil {
		c.JSON(404, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(200, gin.H{
		"id":     id,
		"nombre": nombre,
		"ciudad": ciudad,
		"email":  email,
	})
}
