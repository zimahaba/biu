package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BiuHandler struct {
	DB *gorm.DB
}

func (h BiuHandler) GetUserHandler(c echo.Context) error {
	return c.String(http.StatusOK, "get users "+c.Param("id"))
}

func (h BiuHandler) CreateUserHandler(c echo.Context) error {
	userRequest := UserRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&userRequest)
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userRequest.Password = string(password)

	u := userRequest.ToAppUser()
	result := h.DB.Clauses(clause.Returning{}).Create(&u)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, IdResource{Id: u.ID})
}

func (h BiuHandler) LoginHandler(c echo.Context) error {
	var creds CredentialsRequest
	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		return err
	}

	var passwordHash string
	result := h.DB.Select("password").First(&passwordHash, "username = ?", creds.Username)
	if result.Error != nil {
		return result.Error
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password))
	if err != nil {
		return err
	}

	tokenCookie, err := GenerateTokenCookie(creds.Username)
	if err != nil {
		return err
	}
	c.SetCookie(tokenCookie)

	if creds.KeepLoggedIn {
		refreshToken, refreshCookie, err := GenerateRefreshCookie(creds.Username)
		if err != nil {
			return err
		}

		//err = UpsertRefreshToken(refreshToken, username, db)
		if refreshToken != "nil" {
			return err
		}

		c.SetCookie(refreshCookie)
	}

	return c.NoContent(http.StatusOK)
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
