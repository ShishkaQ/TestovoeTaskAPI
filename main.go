package main

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"go-redis-task-service/handlers"
	"go-redis-task-service/storage"
	"go-redis-task-service/tasks"
)

// @title Task Service API
// @version 1.0
// @description API для управления задачами с использованием Redis и SQLite
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()

	// Инициализация БД
	db, err := storage.InitDB(os.Getenv("SQLITE_PATH"))
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}

	// Инициализация Redis
	redisClient := storage.InitRedis(ctx, &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	// Запуск пула воркеров
	const workerCount = 5
	for i := 1; i <= workerCount; i++ {
		go tasks.Worker(ctx, db, redisClient, i)
	}

	// Роутеры
	router := handlers.NewRouter(db, redisClient)
	log.Println("Starting HTTP server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Gin run error: %v", err)
	}
}