package publisher

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/fikri-yusrihan/outbox-project/internal/model"
	"github.com/segmentio/kafka-go"
)

func StartKafkaPublisher(ctx context.Context, db *sql.DB, id string) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "user.created",
		Balancer: &kafka.LeastBytes{},
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Publisher %s shutting down...", id)
			return

		case <-ticker.C:
			log.Printf("Publisher %s is running...", id)
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				log.Printf("Error starting transaction: %v", err)
				continue
			}

			events, err := model.FetchAndLockUnpublishedEvents(tx, 100)
			if err != nil {
				log.Printf("Error fetching unpublished events: %v", err)
				tx.Rollback()
				continue
			}

			success := true
			for _, event := range events {
				select {
				case <-ctx.Done():
					log.Printf("Publisher %s interrupted during batch", id)
					_ = tx.Rollback()
					return
				default:
				}

				msg, _ := json.Marshal(event.Payload)

				err := writer.WriteMessages(ctx, kafka.Message{
					Key:   []byte(event.EventType),
					Value: msg,
				})

				if err != nil {
					log.Printf("Error publishing event %s: %v", event.EventType, err)
					success = false
					break // stop this batch and rollback
				}

				err = model.MarkEventAsPublished(tx, event.ID)
				if err != nil {
					log.Printf("Error marking event %s as published: %v", event.ID, err)
					success = false
					break
				}
			}

			if success {
				if err := tx.Commit(); err != nil {
					log.Printf("Error committing transaction: %v", err)
				} else {
					log.Printf("Publisher %s Successfully published %d event(s)", id, len(events))
				}
			} else {
				tx.Rollback()
			}
		}
	}
}
