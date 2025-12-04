package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.GET("/cambioEstatus", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cambioEstatus.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.GET("/inicioCliente", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "inicioCliente.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.GET("/listaEnvios", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "listaEnvios.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.GET("/registroEnvio", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "registroEnvio.html", gin.H{
			"title":   "Hello, World!",
			"message": "Welcome to my blog!",
		})
	})

	r.Run(":8080")
}
