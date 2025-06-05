package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartKafkaListener(brokerAddr, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		GroupID:  "user-event-listener",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Println("Kafka listener started for topic:", topic)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		// Process the message
		// Here you would typically handle the event, e.g., update a database or trigger some action
		// For demonstration, we just log the message

		log.Printf("Received message: key=%s value=%s\n", string(msg.Key), string(msg.Value))
	}
}