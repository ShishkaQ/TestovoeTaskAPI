package tasks

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"go-redis-task-service/models"
)

const redisQueueKey = "task_queue"

// Worker читает ID из Redis и обрабатывает задачу
func Worker(ctx context.Context, db *gorm.DB, redisClient *redis.Client, workerID int) {
	for {
		res, err := redisClient.BLPop(ctx, 0, redisQueueKey).Result()
		if err != nil {
			log.Printf("[worker %d] redis BLPop error: %v", workerID, err)
			time.Sleep(time.Second)
			continue
		}
		if len(res) < 2 {
			continue
		}
		var taskID uint
		_ = json.Unmarshal([]byte(res[1]), &taskID)

		log.Printf("[worker %d] picked task %d", workerID, taskID)
		if err := ProcessTask(ctx, db, taskID); err != nil {
			log.Printf("[worker %d] error processing task %d: %v", workerID, taskID, err)
		}
	}
}

// ProcessTask обновляет статус задачи в несколько этапов
func ProcessTask(ctx context.Context, db *gorm.DB, id uint) error {
	_ = updateStatus(db, id, "processing")
	select {
	case <-time.After(10 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	_ = updateStatus(db, id, "checking")

	select {
	case <-time.After(5 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	return updateStatus(db, id, "approved")
}

// updateStatus сохраняет новый статус и инкремент попыток
func updateStatus(db *gorm.DB, id uint, status string) error {
	return db.Model(&models.Task{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "attempts": gorm.Expr("attempts + 1")} ).Error
}
