package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	calc "go_final_project/pkg/calculate"
	chk "go_final_project/pkg/checker"
	"go_final_project/pkg/models"
	"go_final_project/pkg/normilize"
	"go_final_project/pkg/storage"
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
        sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	id, err := storage.AddTaskStorage(h.DB, task)
    if err != nil {
        sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }
	
	sendJSONResponse(w, map[string]interface{}{"id": strconv.Itoa(id)}, http.StatusOK)
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
				sender := map[string]string{"error": err.Error()}
				sendByte, _ := json.Marshal(sender)
				w.Write([]byte(sendByte))
				return
			}
			result["tasks"] = tasks
			sendByte, err := json.Marshal(result)
			if err != nil {
				sender := map[string]string{"error": err.Error()}
				sendByte, _ := json.Marshal(sender)
				w.Write([]byte(sendByte))
				return
			}
			w.Write([]byte(sendByte))
			return
		}
		tasks, err :=storage.SearchTaskToWord(h.DB, search)
		if err != nil {
			sender := map[string]string{"error": err.Error()}
			sendByte, _ := json.Marshal(sender)
			w.Write([]byte(sendByte))
			return
		}
		result["tasks"] = tasks
		sendByte, err := json.Marshal(result)
		if err != nil {
			sender := map[string]string{"error": err.Error()}
			sendByte, _ := json.Marshal(sender)
			w.Write([]byte(sendByte))
			return
		}
		w.Write([]byte(sendByte))
		return
	}

	tasks, err := storage.GetAllTasks(h.DB)
	if err != nil {
		sender := map[string]string{"error": err.Error()}
		sendByte, _ := json.Marshal(sender)
		w.Write([]byte(sendByte))
		return
	}
	
	result["tasks"] = tasks
	sendByte, err := json.Marshal(result)
	if err != nil {
		sender := map[string]string{"error": err.Error()}
		sendByte, _ := json.Marshal(sender)
		w.Write([]byte(sendByte))
		return
	}

	w.Write([]byte(sendByte))
}

func (h handler) GetTaskID(w http.ResponseWriter, r *http.Request) {
    sender := make(map[string]interface{})

    idParam := r.FormValue("id")

    if idParam == "" {
        sender["error"] = "Не указан идентификатор"
        sendJSONResponse(w, sender, http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idParam)
    if err != nil {
        sender["error"] = "Неверный формат идентификатора"
        sendJSONResponse(w, sender, http.StatusBadRequest)
        return
    }

    task, err := storage.GetTaskByID(h.DB, id)
    if err != nil {
        sender["error"] = "Задача не найдена"
        sendJSONResponse(w, sender, http.StatusNotFound)
        return
    }
	byteSend, err := json.Marshal(task)
	if err != nil {
		fmt.Println(task, err)
	}
	w.Write(byteSend)
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
        sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	_, err = storage.GetTaskByID(h.DB, id)
    if err != nil {
        sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	task, err = chk.Task(task)
    if err != nil {
        sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

	if err := storage.UpdateTask(h.DB, task); err != nil {
		sendErrorResponse(w, "Error update on database "+err.Error(), http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, map[string]interface{}{}, http.StatusOK)
}


func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{"error": message})
}


func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(statusCode)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println()
        return
    }
}
