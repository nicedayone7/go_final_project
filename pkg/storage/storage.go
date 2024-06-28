package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go_final_project/pkg/models"

	_ "modernc.org/sqlite"
)

type SchedulerStorage struct {
	db *sql.DB
}

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
	fmt.Println(create)
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
	results, err := db.Query("SELECT * FROM scheduler;")

	if err != nil {
		return nil, err
	}

	var tasks = make([]models.Task, 0)
	for results.Next() {
		i := 0
		var task models.Task
		err = results.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		if i > 10 {
			tasks = append(tasks, task)
			break
		}
		tasks = append(tasks, task)
		i += 1
	}

	return tasks, nil
}
