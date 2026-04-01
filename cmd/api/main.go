package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Divyamsirswal/surge/internal/database"
	"github.com/Divyamsirswal/surge/internal/handlers"
	"github.com/Divyamsirswal/surge/internal/worker"
)

func main() {
	fmt.Println("Booting up Flash Engine Sale !!!")

	database.InitializePostgres()
	database.InitializeRedis()
	database.InitializeMessageQueue()

	defer database.DB.Close()
	defer database.RedisClient.Close()
	defer database.MQConn.Close()
	defer database.MQChannel.Close()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/buy", handlers.FlashSaleHandler)

	port := ":8080"
	fmt.Printf("Server Listening on port %s\n", port)

	go handlers.StartMQPublisher()
	go worker.StartOrderProcessor()

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Server Crashed : ", err)
	}
}
