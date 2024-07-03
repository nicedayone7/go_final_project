package handlers

import "go_final_project/pkg/storage"

type handler struct {
	db storage.Storage
}

func New(db *storage.Storage) handler {
	return handler{
		db: *db,
	}
}