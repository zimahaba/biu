package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppUser struct {
	gorm.Model
	Name  string
	Email string
}

func main() {
	dsn := "host=localhost user=biu password=biu dbname=biu port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("fail")
	}

	e := echo.New()

	h := BiuHandler{db: db}
	e.GET("/users/:id", h.GetUserHandler)
	e.POST("/users", h.CreateUserHandler)
	e.POST("/auth/login", h.LoginHandler)
	e.POST("/auth/logout", h.LogoutHandler)
	e.POST("/auth/refresh", h.RefreshHandler)
	e.POST("/auth/forgot", h.ForgotHandler)
	e.GET("/auth/recover", h.RecoverHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

type BiuHandler struct {
	db *gorm.DB
}

func (h BiuHandler) GetUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, "get users "+c.Param("id"))
}

func (h BiuHandler) CreateUserHandler(c echo.Context) error {
	u := AppUser{Name: "porra", Email: "Carai"}
	h.db.Clauses(clause.Returning{}).Create(&u)
	fmt.Printf("user: %v.\n", u)

	return c.String(http.StatusOK, "create users "+c.Param("id"))
}

func (h BiuHandler) LoginHandler(c echo.Context) error {
	return c.String(http.StatusOK, "login "+c.Param("id"))
}

func (h BiuHandler) LogoutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "logout "+c.Param("id"))
}

func (h BiuHandler) RefreshHandler(c echo.Context) error {
	return c.String(http.StatusOK, "refresh "+c.Param("id"))
}

func (h BiuHandler) ForgotHandler(c echo.Context) error {
	return c.String(http.StatusOK, "forgot "+c.Param("id"))
}

func (h BiuHandler) RecoverHandler(c echo.Context) error {
	return c.String(http.StatusOK, "recover "+c.Param("id"))
}
