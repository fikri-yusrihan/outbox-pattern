package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/fikri-yusrihan/outbox-project/internal/model"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
			return
		}

		defer tx.Rollback() // Ensure rollback on error

		user := model.User{
			Name:  req.Name,
			Email: req.Email,
		}

		err = model.InsertUser(tx, &user)
		if err != nil {
			http.Error(w, "Failed to insert user: ", http.StatusInternalServerError)
			return
		}

		event := model.NewUserCreatedEvent(user)

		err = model.InsertEvent(tx, event)
		if err != nil {
			http.Error(w, "Failed to insert event", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
