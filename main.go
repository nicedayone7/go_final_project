package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)


func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func getEnv(envFile string) string {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Dont load env: %S", err)
	}

	todoPort := os.Getenv("TODO_PORT")
	return todoPort
}

func main() {
	
	r := chi.NewRouter()
	
	fs := http.FileServer(http.Dir("web"))
	r.Handle("/web/*", http.StripPrefix("/web/", fs))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)
	
	port := getEnv(".env")
	http.ListenAndServe(":" + port, r)
}