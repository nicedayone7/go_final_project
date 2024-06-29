package handlers

import (
	"bytes"
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
	var buf bytes.Buffer

	sender := map[string]string{
		"id": "",
		"error": "",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} 

	task, err = chk.Task(task)
	 

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sender["error"] = err.Error()
		sendByte, _ := json.Marshal(sender)
		w.Write(sendByte)
		return
	}
		
	id, err := storage.AddTaskStorage(h.DB, task)
	if err != nil {
		sender["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		sendByte, _ := json.Marshal(sender)
		w.Write(sendByte)
		return
	}

	sender["id"] = fmt.Sprintf("%d", id)
	sendByte, _ := json.Marshal(sender)
	w.WriteHeader(http.StatusOK)
	w.Write(sendByte)
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
				w.Write(sendByte)
				return
			}
			result["tasks"] = tasks
			sendByte, err := json.Marshal(result)
			if err != nil {
				sender := map[string]string{"error": err.Error()}
				sendByte, _ := json.Marshal(sender)
				w.Write(sendByte)
				return
			}
			w.Write(sendByte)
			return
		}
		tasks, err :=storage.SearchTaskToWord(h.DB, search)
		if err != nil {
			sender := map[string]string{"error": err.Error()}
			sendByte, _ := json.Marshal(sender)
			w.Write(sendByte)
			return
		}
		result["tasks"] = tasks
		sendByte, err := json.Marshal(result)
		if err != nil {
			sender := map[string]string{"error": err.Error()}
			sendByte, _ := json.Marshal(sender)
			w.Write(sendByte)
			return
		}
		w.Write(sendByte)
		return
	}

	tasks, err := storage.GetAllTasks(h.DB)
	if err != nil {
		sender := map[string]string{"error": err.Error()}
		sendByte, _ := json.Marshal(sender)
		w.Write(sendByte)
		return
	}
	
	result["tasks"] = tasks
	sendByte, err := json.Marshal(result)
	if err != nil {
		sender := map[string]string{"error": err.Error()}
		sendByte, _ := json.Marshal(sender)
		w.Write(sendByte)
		return
	}

	w.Write(sendByte)
}

func (h handler) GetTaskID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	
	if idParam == "" {
        http.Error(w, "Missing id parameter", http.StatusBadRequest)
        return
    }

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	task, err := storage.GetTaskByID(h.DB ,id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendByte, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(sendByte)
	
}