package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=biu password=biu dbname=biu port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil && db != nil {
		log.Fatal("fail")
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "biu",            // No password set
		DB:       0,                // Use default DB
	})

	e := echo.New()

	h := BiuHandler{DB: db, RC: rc}

	e.GET("/users/:id", h.GetUserHandler)
	e.POST("/users", h.CreateUserHandler)
	e.POST("/auth/login", h.LoginHandler)
	e.POST("/auth/logout", h.LogoutHandler)
	e.POST("/auth/refresh", h.RefreshHandler)
	e.POST("/auth/forgot", h.ForgotHandler)
	e.GET("/auth/recover", h.RecoverHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
