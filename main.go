package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go_final_project/pkg/handlers"
	"go_final_project/pkg/storage"

	"github.com/go-chi/chi/v5"

	// "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	dateFormat      = "20060102"
	storageFilename = "scheduler.db"
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

func getEnv(envFile string) (string, string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Dont load env: %s", err)
	}

	todoPort := os.Getenv("TODO_PORT")
	dbPath := os.Getenv("TODO_DBFILE")
	return todoPort, dbPath
}

func handleRequests(DB *sql.DB, port, workDir string) {
	h := handlers.New(DB)
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web"))
	r.Get("/api/nextdate", h.RequestNextDate)
	r.Post("/api/task", h.AddTask)
	r.Get("/api/tasks", h.GetTasks)
	r.Get("/api/task", h.GetTaskID)
	r.Put("/api/task", h.PutTask)
	r.Post("/api/task/done", h.TaskDone)
	r.Delete("/api/task", h.DeleteTask)
	r.Handle("/web/*", http.StripPrefix("/web/", fs))
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)

	http.ListenAndServe(":" + port, r)
}

func main() {
	port, dbPath := getEnv(".env")
	workDir, _ := os.Getwd()
	fmt.Println(dbPath, workDir)
	if !storage.ExistingStorage(dbPath) {
		if err := storage.CreateStorage(dbPath); err != nil {
			log.Fatalf("Dont create db: %s", err)
		}
		if err := storage.CreateTable(dbPath); err != nil {
			log.Fatalf("Dont create table: %s", err)
		}
	}

	db, err := storage.Connect(dbPath)
	if err != nil {
		log.Fatalf("Dont connect database: %s", err)
	}

	handleRequests(db, port, workDir)
	// r.Use(middleware.Logger)
}
