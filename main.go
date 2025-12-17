package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

type EnvioView struct {
	ID             int
	TrackingNumber string
	Productos      string
	Estatus        string
	Cliente        string
	Fecha          string
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
		rows, err := db.Query("SELECT id, tracking_number, productos, cliente, estatus, fecha_hora FROM envios ORDER BY fecha_hora DESC")
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Error al consultar los envíos")
			log.Println(err)
			return
		}
		defer rows.Close()

		var envios []EnvioView
		for rows.Next() {
			var e Envio
			var eV EnvioView
			if err := rows.Scan(&e.ID, &e.TrackingNumber, &e.Productos, &e.Cliente, &e.Estatus, &e.FechaHora); err != nil {
				ctx.String(http.StatusInternalServerError, "Error al escanear los envíos")
				log.Println(err)
				return
			}
			eV.ID = e.ID
			eV.TrackingNumber = e.TrackingNumber
			eV.Productos = e.Productos
			eV.Cliente = e.Cliente
			eV.Estatus = e.Estatus
			eV.Fecha = e.FechaHora.Format("2006-01-02 15:04")
			envios = append(envios, eV)
		}

		ctx.HTML(http.StatusOK, "listaEnvios.html", gin.H{
			"envios": envios,
		})
	})

	r.GET("/registroEnvio", func(ctx *gin.Context) {
		trackingSKU := ctx.Query("trackingSKU")
		ctx.HTML(http.StatusOK, "registroEnvio.html", gin.H{
			"title":       "Registrar Envío",
			"message":     "Welcome to my blog!",
			"TrackingSKU": trackingSKU,
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
		restaurante := c.PostForm("restaurante")
		producto := c.PostForm("producto")
		estatus := "enviado"
		ahora := time.Now()

		var trackingSKU string = strconv.Itoa(rand.Intn(1000000))
		fmt.Print(trackingSKU)
		sqlStatement := `
			INSERT INTO envios (productos, cliente, tracking_number, estatus, fecha_hora)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`

		var newID int

		err := db.QueryRow(sqlStatement, producto, restaurante, trackingSKU, estatus, ahora).Scan(&newID)
		if err != nil {
			fmt.Println(err)
			c.Redirect(http.StatusFound, "/registroEnvio")
		} else {
			fmt.Println("ID del nuevo envío:", newID)
			fmt.Print("Numero de envio", trackingSKU)
			c.Redirect(http.StatusFound, "/registroEnvio?trackingSKU="+trackingSKU)
		}

	})

	r.POST("/cambioEstatus", func(ctx *gin.Context) {
		trackingNumber := ctx.PostForm("id")
		estatus := ctx.PostForm("status")
		fmt.Println(trackingNumber)
		// ctx.Redirect(http.StatusFound, "/cambioEstatus")
		sqlStatement := `
			UPDATE envios
			SET estatus = $1
			WHERE tracking_number = $2`

		_, err := db.Exec(sqlStatement, estatus, trackingNumber)
		if err != nil {
			fmt.Println(err)
			ctx.Redirect(http.StatusFound, "/cambioEstatus")
		} else {
			fmt.Println("Estatus actualizado")
			ctx.Redirect(http.StatusFound, "/cambioEstatus")
		}
	})

	r.POST("/obtenerEnvio", func(ctx *gin.Context) {
		trackingNumberUser := ctx.PostForm("tracking_number")
		fmt.Println("Numero de envio", trackingNumberUser)
		sqlStatement := `
			SELECT productos, cliente, tracking_number, estatus, fecha_hora
			FROM envios
			WHERE tracking_number = $1`
		rows, err := db.Query(sqlStatement, trackingNumberUser)
		if err != nil {
			fmt.Println(err)
			ctx.Redirect(http.StatusFound, "/inicioCliente")
		}
		defer rows.Close()

		var productos, cliente, trackingNumber, estatus string
		var fechaHora time.Time

		if rows.Next() {
			err := rows.Scan(&productos, &cliente, &trackingNumber, &estatus, &fechaHora)
			if err != nil {
				fmt.Println(err)
				ctx.HTML(http.StatusOK, "inicioCliente.html", gin.H{
					"resultado": "Error al procesar el envío.",
					"color":     "hsl(348, 100%, 61%)",
					"guia":      trackingNumberUser,
				})
				return
			}
			ctx.HTML(http.StatusOK, "inicioCliente.html", gin.H{
				"encontrado": true,
				"productos":  productos,
				"cliente":    cliente,
				"guia":       trackingNumber,
				"estatus":    estatus,
				"fechaHora":  fechaHora.Format("02/01/2006 15:04"),
			})
		} else {
			ctx.HTML(http.StatusOK, "inicioCliente.html", gin.H{
				"resultado": "No se encontró ningún envío con ese número de guía.",
				"color":     "hsl(48, 100%, 67%)",
				"guia":      trackingNumberUser,
			})
		}
	})

	r.Run(":8080")
}
