package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Event struct {
	ID          string     `json:"id"`
	EventType   string     `json:"event_type"`
	Payload     string     `json:"payload"`
	Published   bool       `json:"published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func NewUserCreatedEvent(user User) Event {
	payload, _ := json.Marshal(user)

	return Event{
		EventType: "user.created",
		Payload:   string(payload),
	}
}

func InsertEvent(tx *sql.Tx, event Event) error {
	query := `
		INSERT INTO events (event_type, payload)
		VALUES ($1, $2)
	`
	_, err := tx.Exec(query, event.EventType, event.Payload)
	return err
}

func FetchAndLockUnpublishedEvents(tx *sql.Tx, limit int) ([]Event, error) {
	query := `
		SELECT id, event_type, payload, published, published_at, created_at
		FROM events
		WHERE published = false
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT $1
	`
	rows, err := tx.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var publishedAt sql.NullTime

		if err := rows.Scan(&event.ID, &event.EventType, &event.Payload, &event.Published, &publishedAt, &event.CreatedAt); err != nil {
			return nil, err
		}

		event.PublishedAt = &publishedAt.Time
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func MarkEventAsPublished(tx *sql.Tx, eventID string) error {
	query := `
		UPDATE events
		SET published = true, published_at = NOW()
		WHERE id = $1
	`
	_, err := tx.Exec(query, eventID)
	return err
}