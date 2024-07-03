package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go_final_project/pkg/auth"
	"go_final_project/pkg/models"
	"go_final_project/pkg/nextdate"
	"go_final_project/pkg/normilize"
	"go_final_project/pkg/storage"
	"go_final_project/pkg/wrapper"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func HandleRequests(dbPath, port string) {
	if err := storage.StartStorage(dbPath); err != nil {
		log.Fatalf("Error start storage: %s", err)
	}
	db, err := storage.Connect(dbPath)
	if err != nil {
		log.Fatalf("Error connection storage: %s", err)
	}
	stor := storage.New(db)
	h := New(&stor)
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/api/signin", h.Auth)
	r.Get("/api/nextdate", h.RequestNextDate)
	r.With(auth.AuthMiddleware).Post("/api/task", h.AddTask)
	r.With(auth.AuthMiddleware).Get("/api/tasks", h.GetTasks)
	r.With(auth.AuthMiddleware).Get("/api/task", h.GetTaskID)
	r.With(auth.AuthMiddleware).Put("/api/task", h.PutTask)
	r.With(auth.AuthMiddleware).Post("/api/task/done", h.TaskDone)
	r.With(auth.AuthMiddleware).Delete("/api/task", h.DeleteTask)
	r.Get("/*", http.FileServer(http.Dir("web")).ServeHTTP)

	log.Printf("server start on port %s", port)
	http.ListenAndServe(":"+port, r)
}

func (h handler) RequestNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	startDate := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse(nextdate.DateFormat, now)
	if err != nil {
		fmt.Fprintf(w, "Error parse date")
	}

	nextDateTask, err := nextdate.NextDate(nowTime, startDate, repeat)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	w.Write([]byte(nextDateTask))
}

func (h handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	task, err := task.Check()
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.db.AddTaskStorage(task)
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	wrapper.SendJSONResponse(w, map[string]interface{}{"id": strconv.Itoa(id)}, http.StatusOK)
}

func (h handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	result := make(map[string][]models.Task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if search != "" {
		dateSearch := normilize.SearchTasks(search)

		if dateSearch != "" {
			tasks, err := h.db.SearchTaskToDate(dateSearch)
			if err != nil {
				wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}

			result["tasks"] = tasks
			wrapper.SendJSONResponse(w, result, http.StatusOK)
			return
		}

		tasks, err := h.db.SearchTaskToWord(search)
		if err != nil {
			wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		result["tasks"] = tasks
		wrapper.SendJSONResponse(w, result, http.StatusOK)
		return
	}

	tasks, err := h.db.GetAllTasks()
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	result["tasks"] = tasks
	wrapper.SendJSONResponse(w, result, http.StatusOK)
}

func (h handler) GetTaskID(w http.ResponseWriter, r *http.Request) {
	sender := make(map[string]interface{})

	idParam := r.FormValue("id")

	if idParam == "" {
		sender["error"] = "Не указан идентификатор"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		sender["error"] = "Неверный формат идентификатора"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	task, err := h.db.GetTaskByID(id)
	if err != nil {
		sender["error"] = "Задача не найдена"
		wrapper.SendJSONResponse(w, sender, http.StatusNotFound)
		return
	}

	wrapper.SendJSONResponse(w, task, http.StatusOK)
}

func (h handler) PutTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Error decoding JSON "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.db.GetTaskByID(id)
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err = task.Check()
	if err != nil {
		wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateTask(task); err != nil {
		wrapper.SendErrorResponse(w, "Error update on database "+err.Error(), http.StatusBadRequest)
		return
	}

	wrapper.SendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}

func (h handler) TaskDone(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	sender := make(map[string]interface{})
	idParam := r.FormValue("id")

	if idParam == "" {
		sender["error"] = "Не указан идентификатор"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		sender["error"] = "Неверный формат идентификатора"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	task, err = h.db.GetTaskByID(id)
	if err != nil {
		sender["error"] = "Задача не найдена"
		wrapper.SendJSONResponse(w, sender, http.StatusNotFound)
		return
	}

	if task.Repeat != "" {
		task.Date, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			wrapper.SendErrorResponse(w, "Error calculate next date "+err.Error(), http.StatusBadRequest)
			return
		}
		err = h.db.UpdateTask(task)
		if err != nil {
			wrapper.SendErrorResponse(w, "Error update task in database "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if err := h.db.DeleteTask(id); err != nil {
			wrapper.SendErrorResponse(w, "Error delete task in database "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	wrapper.SendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}

func (h handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	sender := make(map[string]interface{})
	idParam := r.FormValue("id")

	if idParam == "" {
		sender["error"] = "Не указан идентификатор"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		sender["error"] = "Неверный формат идентификатора"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}

	if err = h.db.DeleteTask(id); err != nil {
		sender["error"] = "Error delete DB"
		wrapper.SendJSONResponse(w, sender, http.StatusBadRequest)
		return
	}
	wrapper.SendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}

func (h handler) Auth(w http.ResponseWriter, r *http.Request) {
	var pass models.Password
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&pass); err != nil {
		http.Error(w, "Error decoding JSON "+err.Error(), http.StatusBadRequest)
		return
	}
	originalPass := os.Getenv("TODO_PASSWORD")
	if originalPass == pass.SendPassword {
		jwtToken, err := SignedToken()
		if err != nil {
			wrapper.SendJSONResponse(w, map[string]interface{}{"error": "Неверный пароль"}, http.StatusBadRequest)
			return
		}
		wrapper.SendJSONResponse(w, map[string]interface{}{"token": jwtToken}, http.StatusAccepted)
		return
	}

	wrapper.SendJSONResponse(w, map[string]interface{}{"error": "Неверный пароль"}, http.StatusBadRequest)
}

func SignedToken() (string, error) {
	secret := os.Getenv("TODO_SECRET")
	secretStr := []byte(secret)
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secretStr)
	return signedToken, err
}
