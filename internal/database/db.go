package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

var DB *sql.DB
var RedisClient *redis.Client

func InitializePostgres() {
	dsn := "postgres://admin:supersecret@localhost:5432/surge_db?sslmode=disable"
	var err error

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to load Postgres driver : ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Postgres is down or credentials wrong : ", err)
	}

	fmt.Println("Database Connected Successfully.")

}

func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		PoolSize:     10000,
		MinIdleConns: 100,
	})

	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Redis is unreachable : ", err)
	}

	fmt.Println("Redis Connected Successfully.")
}
