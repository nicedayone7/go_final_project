package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	webDir = "../web"
	serverAddr = "localhost:7540"
)

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".webindex.html")
	// http.ServeFile(w, r, "./web/js/")
}

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "web\\index.html")
	})
	// http.Handle("/", http.FileServer(http.Dir(webDir))) 
	http.ListenAndServe(serverAddr, nil)
}