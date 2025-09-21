package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Publication struct {
	Title   string
	Content string
}

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL no está configurada")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error al hacer ping a la base de datos: ", err)
	}

	fmt.Println("Conexión a la base de datos exitosa!")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Configura el manejador para la ruta principal "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Publication{
			Title:   "Mi primer publicación",
			Content: "Contenido de mi primer publicación, ¡generado con Go!",
		}

		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			log.Printf("Error al parsear la plantilla: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error al ejecutar la plantilla: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	})

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
