package storage

import (
	"database/sql"
	"fmt"
	"go_final_project/pkg/models"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	storageFilename        = "scheduler.db"
	limit = "10"
	create          string = `
CREATE TABLE IF NOT EXISTS scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date VARCHAR(256) NOT NULL,
title VARCHAR(256) NOT NULL,
comment VARCHAR(1024) NOT NULL,
repeat VARCHAR(128) NOT NULL
);`
)

func StartStorage(dbPath string) error {
	if !ExistingStorage(dbPath) {
		if err := CreateStorage(dbPath); err != nil {
			log.Fatalf("Error create db: %s", err)
		}
		if err := CreateTable(dbPath); err != nil {
			log.Fatalf("Error create table: %s", err)
		}
	}
	return nil
}

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

func CreateTable(storagePath string) error {
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		log.Fatalf("Dont connect to database: %s", err)
	}

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

	err = db.Ping()
    if err != nil {
        return nil, err
    }

	return db, err
} 

func (s *Storage) AddTaskStorage(task models.Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
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
	fmt.Println(id)
	return int(id), nil
}

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	rows, err := s.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit;", sql.Named("limit", limit))

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

func (s *Storage) SearchTaskToWord(search string) ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
    rows, err := s.db.Query(query, sql.Named("search", "%"+search+"%"), sql.Named("limit", limit))
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) SearchTaskToDate(dateToSearch string) ([]models.Task, error) {
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit", sql.Named("date", dateToSearch), sql.Named("limit", limit))
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) GetTaskByID(id int) (models.Task, error) {
	var task models.Task 
	row := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	_, err := s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
	sql.Named("date", task.Date),
	sql.Named("title", task.Title),
	sql.Named("comment", task.Comment),
	sql.Named("repeat", task.Repeat),
	sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	
	return nil
}