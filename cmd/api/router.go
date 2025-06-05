package main

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", CreateUserHandler(db))

	return mux
}