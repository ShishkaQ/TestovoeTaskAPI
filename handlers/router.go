package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-redis-task-service/docs" // Импорт docs
	"go-redis-task-service/handlers/task"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func NewRouter(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	r := gin.Default()
	
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	
	// API routes
	r.POST("/task", task.CreateTask(db, redisClient))
	r.GET("/task/:id", task.GetTaskByID(db))
	
	// Инициализация Swagger info
	docs.SwaggerInfo.Title = "Task Service API"
	docs.SwaggerInfo.Description = "API для управления задачами"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	return r
}