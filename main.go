package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"stu-go/common"
)

// Config holds runtime configuration values.
type Config struct {
	Port         string
	MySQLDSN     string
	RedisAddr    string
	RedisPass    string
	RedisDB      int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// App aggregates shared dependencies for handlers.
type App struct {
	DB    *sql.DB
	Redis *redis.Client
}

func main() {
	cfg := loadConfig()

	db, err := common.OpenMySQL(common.MySQLConfig{DSN: cfg.MySQLDSN})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	redisClient, err := common.OpenRedis(common.RedisConfig{Addr: cfg.RedisAddr, Password: cfg.RedisPass, DB: cfg.RedisDB})
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	app := &App{DB: db, Redis: redisClient}
	router := setupRouter(app)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Printf("server listening on :%s", cfg.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}

func loadConfig() Config {
	return Config{
		Port:         getEnv("PORT", "8080"),
		MySQLDSN:     getEnv("MYSQL_DSN", "user:password@tcp(localhost:3306)/dbname?parseTime=true"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:      0,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func setupRouter(app *App) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		mysqlStatus := "ok"
		redisStatus := "ok"

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := app.DB.PingContext(ctx); err != nil {
			mysqlStatus = err.Error()
		}
		if err := app.Redis.Ping(ctx).Err(); err != nil {
			redisStatus = err.Error()
		}

		c.JSON(http.StatusOK, gin.H{
			"mysql": mysqlStatus,
			"redis": redisStatus,
			"time":  time.Now().UTC(),
		})
	})

	router.GET("/now", func(c *gin.Context) {
		var current time.Time
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := app.DB.QueryRowContext(ctx, "SELECT NOW()").Scan(&current); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"mysql_time": current})
	})

	router.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		val, err := app.Redis.Get(ctx, key).Result()
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
	})

	router.POST("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		var payload struct {
			Value      string        `json:"value" binding:"required"`
			Expiration time.Duration `json:"expiration"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := app.Redis.Set(ctx, key, payload.Value, payload.Expiration).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "stored"})
	})

	return router
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
