package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Envio struct {
	ID             int
	Productos      string
	Cliente        string
	TrackingNumber string
	Estatus        string
	FechaHora      time.Time
}

func main() {
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Error conectando a la BD: ", err)
	}

	queryCrearTabla := `CREATE TABLE IF NOT EXISTS envios (
		id SERIAL PRIMARY KEY,
		productos TEXT NOT NULL,
		cliente VARCHAR(255) NOT NULL,
		tracking_number VARCHAR(100) UNIQUE NOT NULL,
		estatus VARCHAR(50) DEFAULT 'creado',
		fecha_hora TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(queryCrearTabla); err != nil {
		log.Fatal("Error creando la tabla: ", err)
	}
	fmt.Println("Base de datos verificada correctamente.")
	rows, err := db.Query("SELECT id, productos, cliente, tracking_number, estatus, fecha_hora FROM envios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tCLIENTE\tTRACKING\tESTATUS\tFECHA\tPRODUCTOS")
	fmt.Fprintln(w, "--\t-------\t--------\t-------\t-----\t---------")

	for rows.Next() {
		var e Envio
		err := rows.Scan(&e.ID, &e.Productos, &e.Cliente, &e.TrackingNumber, &e.Estatus, &e.FechaHora)
		if err != nil {
			log.Fatal(err)
		}

		fechaFormateada := e.FechaHora.Format("2006-01-02 15:04")

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			e.ID, e.Cliente, e.TrackingNumber, e.Estatus, fechaFormateada, e.Productos)
	}

	w.Flush()

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

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

	r.POST("/login", func(c *gin.Context) {
		usuario := c.PostForm("username")
		password := c.PostForm("password")

		if usuario == "Almacen" && password == "Almacen" {
			fmt.Print("Contraseña correcta")
			c.Redirect(http.StatusFound, "/cambioEstatus")

		} else if usuario == "Secretaria" && password == "Secretaria" {
			fmt.Print("Contraseña correcta")
			c.Redirect(http.StatusFound, "/registroEnvio")
		} else {
			fmt.Print("Contraseña incorrecta")
			fmt.Println(usuario)
			fmt.Println(password)
			c.Redirect(http.StatusFound, "/login")
		}
	})

	r.POST("/crearEnvio", func(c *gin.Context) {
		restaurante := c.PostForm("rest")
		producto := c.PostForm("prod")

		fmt.Println(restaurante)
		fmt.Println(producto)
		c.Redirect(http.StatusFound, "/registroEnvio")
	})

	r.Run(":8080")
}
