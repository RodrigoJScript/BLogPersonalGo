package main

import (
	"html/template"
	"log"
	"net/http"
	// "database/sql"
	_ "github.com/lib/pq"
)

const (
	validUsername = "Squad"
	validPassword = "Squad2020"
)

func main() {
	// Descomenta este bloque al integrar la base de datos
	/*
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
	*/

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", serveLogin)
	http.HandleFunc("/login", processLogin)
	http.HandleFunc("/dashboard", serveDashboard)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
	tmpl.Execute(w, nil)
}

func processLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == validUsername && password == validPassword {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		log.Println("Credenciales válidas")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Credenciales inválidas")
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title   string
		Content string
	}{
		Title:   "Mi primer publicación",
		Content: "Contenido de mi primer publicación, ¡generado con Go!",
	}

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, data)
}
