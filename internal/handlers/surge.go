package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Divyamsirswal/surge/internal/database"
	"github.com/Divyamsirswal/surge/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

var localOrderQueue = make(chan models.OrderMessage, 50000)

func FlashSaleHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req models.BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Success: false,
			Message: "Invalid Request Payload",
		})
		return
	}

	redisKey := fmt.Sprintf("product:%d:stock", req.ProductID)
	ctx := context.Background()

	stockLeft, err := database.RedisClient.Decr(ctx, redisKey).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Success: false,
			Message: "Cache Error",
		})
		return
	}

	if stockLeft < 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.ApiResponse{
			Success: false,
			Message: "Sold Out! Better luck next time.",
		})
		return
	}

	orderMsg := models.OrderMessage{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Status:    "PENDING_DB_WRITE",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	localOrderQueue <- orderMsg

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ApiResponse{
		Success: true,
		Message: "Order Placed Successfully! You are in the queue.",
	})
}

func StartMQPublisher() {
	for order := range localOrderQueue {
		body, _ := json.Marshal(order)
		database.MQChannel.PublishWithContext(context.Background(),
			"",
			database.OrderQueueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
	}
}
