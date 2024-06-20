package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	calc "go_final_project/calculate"
	st "go_final_project/storage"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	dateFormat = "20060102"
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

func requestNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	fmt.Println(date)
	fmt.Println(now, date,repeat)
	
	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		fmt.Fprintf(w, "Error parse date")
	}

	nextDateTask, err := calc.NextDate(nowTime, date, repeat)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	w.Write([]byte(nextDateTask))
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

func main() {
	
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	
	fs := http.FileServer(http.Dir("web"))
	r.Get("/api/nextdate", requestNextDate)
	r.Handle("/web/*", http.StripPrefix("/web/", fs))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	FileServer(r, "/", filesDir)
	
	port, dbPath := getEnv(".env")

	install := st.ExistingStorage(dbPath)
	if install {
		_, err := st.CreateStorage(dbPath)
		if err != nil {
			log.Fatalf("Dont create db: %s", err) 
		}
	}
	http.ListenAndServe(":" + port, r)

	
}