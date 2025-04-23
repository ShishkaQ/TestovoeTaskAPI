package task

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"go-redis-task-service/models"
)

const redisQueueKey = "task_queue"


// CreateTask godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу и помещает в очередь Redis
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Данные задачи"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]interface{}
// @Router /task [post]
// CreateTask сохраняет задачу и ставит ID в очередь Redis
func CreateTask(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task := models.Task{
			Title:       input.Title,
			Description: input.Description,
			Status:      "new",
			Attempts:    0,
		}
		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create task"})
			return
		}

		payload, _ := json.Marshal(task.ID)
		_ = redisClient.RPush(c, redisQueueKey, payload).Err()

		c.JSON(http.StatusOK, task)
	}
}


// GetTaskByID godoc
// @Summary Получить задачу по ID
// @Description Возвращает задачу из базы данных по её ID
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]interface{}
// @Router /task/{id} [get]
// GetTaskByID возвращает задачу по ID
func GetTaskByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var t models.Task
		if err := db.First(&t, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	}
}