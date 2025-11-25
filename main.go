package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Printf("Error al parsear la plantilla: %v", err)
		http.Error(w, "Error interno del servidor.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error al ejecutar la plantilla: %v", err)
		http.Error(w, "Error al generar la página.", http.StatusInternalServerError)
		return
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("./templates/blog/blog.html")
	if err != nil {
		log.Printf("Error al parsear la plantilla: %v", err)
		http.Error(w, "Error interno del servidor.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error al ejecutar la plantilla: %v", err)
		http.Error(w, "Error al generar la página.", http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Println("Starting server on port 8080...")
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/blog", blogHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
