package model

import "database/sql"

type User struct {
	ID    string  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func InsertUser(tx *sql.Tx, user *User) error {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	err := tx.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	return err
}