package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	calc "go_final_project/pkg/calculate"
	chk "go_final_project/pkg/checker"
	"go_final_project/pkg/models"
	"go_final_project/pkg/normilize"
	"go_final_project/pkg/storage"
	"go_final_project/pkg/wrapper"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

const dateFormat = "20060102"

type handler struct {
	DB *sql.DB
}

func New(db *sql.DB) handler {
	return handler{db}
}

func (h handler) RequestNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

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

func (h handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
	if err := decoder.Decode(&task); err != nil {
        http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

	task, err := chk.Task(task)
    if err != nil {
        wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	id, err := storage.AddTaskStorage(h.DB, task)
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
			tasks, err := storage.SearchTaskToDate(h.DB, dateSearch)
			if err != nil {
				wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}

			result["tasks"] = tasks
			wrapper.SendJSONResponse(w, result, http.StatusOK)
			return
		}

		tasks, err :=storage.SearchTaskToWord(h.DB, search)
		if err != nil {
			wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		result["tasks"] = tasks
		wrapper.SendJSONResponse(w, result, http.StatusOK)
		return
	}

	tasks, err := storage.GetAllTasks(h.DB)
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

    task, err := storage.GetTaskByID(h.DB, id)
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

	_, err = storage.GetTaskByID(h.DB, id)
    if err != nil {
        wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	task, err = chk.Task(task)
    if err != nil {
        wrapper.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	if err := storage.UpdateTask(h.DB, task); err != nil {
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

    task, err = storage.GetTaskByID(h.DB, id)
    if err != nil {
        sender["error"] = "Задача не найдена"
        wrapper.SendJSONResponse(w, sender, http.StatusNotFound)
        return
    }

	if task.Repeat != "" {
		task.Date, err = calc.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			wrapper.SendErrorResponse(w, "Error calculate next date " + err.Error(), http.StatusBadRequest)
			return
		}
		err = storage.UpdateTask(h.DB, task)
		if err != nil {
			wrapper.SendErrorResponse(w, "Error update task in database " + err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if err := storage.DeleteTask(h.DB, id); err != nil {
			wrapper.SendErrorResponse(w, "Error delete task in database " + err.Error(), http.StatusBadRequest)
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

	if err = storage.DeleteTask(h.DB, id); err != nil {
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
	originalPass := getEnvPass(".env")
	if originalPass == pass.SendPassword {
		jwtToken, err := SignedToken(".env")
		if err != nil {
			wrapper.SendJSONResponse(w, map[string]interface{}{ "error": "Неверный пароль"}, http.StatusBadRequest)
			return
		}
		wrapper.SendJSONResponse(w, map[string]interface{}{"token": jwtToken}, http.StatusAccepted)
		return
	}
	
	wrapper.SendJSONResponse(w, map[string]interface{}{ "error": "Неверный пароль"}, http.StatusBadRequest)
}



func SignedToken(pathToEnv string) (string, error) {
	secret := getEnvSecret(pathToEnv)
	secretStr := []byte(secret)
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString(secretStr)
	return signedToken, err
}


func getEnvPass(envFile string) string {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Dont load env pass: %s", err)
	}

	password := os.Getenv("TODO_PASSWORD")
	return password
}

func getEnvSecret(envFile string) string {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Dont load env secret: %s", err)
	}

	password := os.Getenv("TODO_SECRET")
	return password
}