package handlers

import "github.com/gin-gonic/gin"

// Index renderiza la página principal
func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}
