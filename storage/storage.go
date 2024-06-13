package Storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

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
repeat VARCHAR(32) NOT NULL
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

func CreateStorage(storagePath string) (*sql.DB, error) {
	
	_, err := os.Create(filepath.Join(storagePath, storageFilename))
	if err != nil {
		log.Fatalf("Dont create db file: %s", err)
	}
	db, err := sql.Open("sqlite", storageFilename)
	if err != nil {
		log.Fatalf("Dont connect to database: %s", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Hour)

	if _, err = db.Exec(create); err != nil {
		log.Fatalf("Dont create table in database: %s", err)
	}
	return db, nil
}
