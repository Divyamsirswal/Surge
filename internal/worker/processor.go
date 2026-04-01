package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Divyamsirswal/surge/internal/database"
	"github.com/Divyamsirswal/surge/internal/models"
)

func StartOrderProcessor() {
	messages, err := database.MQChannel.Consume(
		database.OrderQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Worker failed to register with RabbitMQ: ", err)
	}

	fmt.Println("Background Worker Started: Listening for orders to save in Postgres...")
	for msg := range messages {
		var order models.OrderMessage
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			log.Println("Failed to parse message, dropping it: ", err)
			msg.Nack(false, false)
			continue
		}

		query := `INSERT INTO orders (user_id, product_id, status) VALUES ($1, $2, $3)`
		_, err = database.DB.ExecContext(context.Background(), query, order.UserID, order.ProductID, "COMPLETED")
		if err != nil {
			log.Println("Database Insert Failed: ", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
		fmt.Printf("DB Saved -> User: %d got Product: %d\n", order.UserID, order.ProductID)
	}
}
