package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/handlers"
	"go_final_project/pkg/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func getEnv(envFile string) (string, string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Dont load env: %s", err)
	}

	todoPort := os.Getenv("TODO_PORT")
	dbPath := os.Getenv("TODO_DBFILE")

	return todoPort, dbPath
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Dont load env: %s", err)
		}
		pass := os.Getenv("TODO_PASSWORD")
		secret := os.Getenv("TODO_SECRET")
		if len(pass) > 0 {
			var jwtStr string
			cookieToken, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookieToken.Value
			}

			token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing token: %v", err), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}

func handleRequests(DB *sql.DB, port string) {
	h := handlers.New(DB)
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	
	r.Post("/api/signin", h.Auth)
	r.Get("/api/nextdate", h.RequestNextDate)
	r.With(authMiddleware).Post("/api/task", h.AddTask)
	r.With(authMiddleware).Get("/api/tasks", h.GetTasks)
	r.With(authMiddleware).Get("/api/task", h.GetTaskID)
	r.With(authMiddleware).Put("/api/task", h.PutTask)
	r.With(authMiddleware).Post("/api/task/done", h.TaskDone)
	r.With(authMiddleware).Delete("/api/task", h.DeleteTask)
	r.Get("/*", http.FileServer(http.Dir("web")).ServeHTTP)

	http.ListenAndServe(":" + port, r)
}

func main() {
	port, dbPath := getEnv(".env")
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

	handleRequests(db, port)
}
