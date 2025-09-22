package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "database/sql"
	_ "github.com/lib/pq"
)

const (
	validUsername = "squad"
	validPassword = "squad2020"
)

var sessions = map[string]string{}

func main() {
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

	http.HandleFunc("/", handleLogin)
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/logout", logout)

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Error bool
	}{
		Error: false,
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == validUsername && password == validPassword {
			sessionToken := uuid.NewString()
			sessions[sessionToken] = username
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: time.Now().Add(120 * time.Minute),
			})
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			log.Println("Credenciales válidas, sesión iniciada")
			return
		}
		data.Error = true
		log.Println("Credenciales inválidas")
	}

	tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
	tmpl.Execute(w, data)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.Error(w, "Error de servidor", http.StatusInternalServerError)
		return
	}
	sessionToken := c.Value

	username, ok := sessions[sessionToken]
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	log.Printf("Usuario '%s' ha accedido a /home", username)
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

func logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	delete(sessions, c.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
