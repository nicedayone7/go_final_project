package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	calc "go_final_project/calculate"
	chk "go_final_project/pkg/checker"
	"go_final_project/pkg/models"
	"go_final_project/pkg/storage"
	"go_final_project/pkg/sorted"
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

	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} 

	task, err = chk.Task(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8") 

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
	tasks, err := storage.GetAllTasks(h.DB)
	fmt.Println(tasks)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		sender := map[string]string{"error": err.Error()}
		sendByte, _ := json.Marshal(sender)
		w.Write(sendByte)
		return
	}
	tasks = wraper.SortedTasks(tasks)

	result := make(map[string][]models.Task)
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