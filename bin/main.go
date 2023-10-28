package main

import (
	"log"
	"net/http"
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ServerConfig struct {
	PORT string
}

func main() {
	var config ServerConfig = ServerConfig{PORT: ":5000"}
	r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	configurePageRoutes(r)

	log.Println("Starting server on port ", config.PORT)
	http.ListenAndServe(config.PORT, r)
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle errs 

		w.WriteHeader(503)
		w.Write([]byte("Bad"))
	}
}


func configurePageRoutes(r *chi.Mux) {
	// these pages are not part of the api, they will be for the 
	// visitor to display the purpose of the image server. 
	// Ex: cr
	r.Method(http.MethodGet, "/", Handler(indexHandler))
}

type Page struct {
	Page string
	Title string
	CSS string
}

func (p Page) renderTemplate(w http.ResponseWriter) {
	// usually i have this as a stand alone func but i think
	// calling it as a method may be more cooler of me
	tmpl, err := template.ParseFiles("templates/" + p.Page)
    if err != nil {
        panic(err)
    }
    err = tmpl.Execute(w, p)
    if err != nil {
        panic(err)
    }
}

var (
	css = "/static/css/main.css"
)


func indexHandler(w http.ResponseWriter, r *http.Request) error {
    p := Page{
        Page: "index.html",
        Title: "chi-img-server",
        CSS: css,
    }
    // this prevents the superflous hanlder err 
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    p.renderTemplate(w)
    return nil
}

// https://go-chi.io/#/pages/routing