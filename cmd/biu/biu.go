package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zimahaba/biu/internal/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	loadEnv()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil && db != nil {
		log.Fatal("fail")
	}
	log.Printf("Connected to db: %v\n", dsn)

	rsn := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	/*rc := redis.NewClient(&redis.Options{
		Addr:     rsn,
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})*/
	log.Printf("Connected to redis: %v\n", rsn)

	e := echo.New()

	h := handlers.BiuHandler{DB: db}

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/users/:id", h.GetUserHandler)
	e.POST("/users", h.CreateUserHandler)
	e.POST("/auth/login", h.LoginHandler)
	e.POST("/auth/logout", h.LogoutHandler)
	e.POST("/auth/refresh", h.RefreshHandler)
	e.POST("/auth/forgot", h.ForgotHandler)
	e.GET("/auth/recover", h.RecoverHandler)
	e.GET("/auth/verify", h.RecoverHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Authorization", "Content-type"},
		AllowCredentials: true,
	}))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
