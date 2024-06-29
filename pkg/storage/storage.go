package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"go_final_project/pkg/models"

	_ "modernc.org/sqlite"
)

const (
	storageFilename        = "scheduler.db"
	create          string = `
CREATE TABLE IF NOT EXISTS scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date VARCHAR(256) NOT NULL,
title VARCHAR(256) NOT NULL,
comment VARCHAR(1024) NOT NULL,
repeat VARCHAR(128) NOT NULL
);`
)

func ExistingStorage(path string) bool {
	storageFile := filepath.Join(filepath.Dir(path), storageFilename)
	_, err := os.Stat(storageFile)

	var install bool
	if err != nil {
		install = true
	}

	return install
}

func CreateStorage(storagePath string) error {
	_, err := os.Create(storagePath)
	if err != nil {
		log.Fatalf("Dont create db file: %s", err)
	}
	return nil
}

func CreateTable(storageFilename string) error {
	db, err := sql.Open("sqlite", storageFilename)
	if err != nil {
		log.Fatalf("Dont connect to database: %s", err)
	}
	defer db.Close()

	if _, err = db.Exec(create); err != nil {
		log.Fatalf("Dont create table in database: %s", err)
	}
	return nil
}

func Connect(storageFilename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", storageFilename)
	if err != nil {
        return nil, err
    }

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
    if err != nil {
        return nil, err
    }

	return db, err
} 

func AddTaskStorage(db *sql.DB,task models.Task) (int, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
						sql.Named("date", task.Date),
						sql.Named("title", task.Title),
						sql.Named("comment", task.Comment),
						sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetAllTasks(db *sql.DB) ([]models.Task, error) {
	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit;", sql.Named("limit", "10"))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks = make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func SearchTaskToWord(db *sql.DB, search string) ([]models.Task, error) {
	query := "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
    rows, err := db.Query(query, sql.Named("search", "%"+search+"%"), sql.Named("limit", "10"))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

	var tasks = make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func SearchTaskToDate(db *sql.DB, dateToSearch string) ([]models.Task, error) {
	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit", sql.Named("date", dateToSearch), sql.Named("limit", "10"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks = make([]models.Task, 0)
	for rows.Next() {
		var task models.Task 
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		} 
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(db *sql.DB, id int) (models.Task, error) {
	var task models.Task 
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}