package main

import (
	"html/template"
	"log"
	"net/http"
)

type Publication struct {
	Title   string
	Content string
}

func main() {

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
