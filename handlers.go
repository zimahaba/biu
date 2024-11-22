package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BiuHandler struct {
	DB *gorm.DB
	RC *redis.Client
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
	result := h.DB.Model(&UserCredentials{}).Select("password").First(&passwordHash, "username = ?", creds.Username)
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

		result = h.DB.Model(&UserCredentials{}).Where("username = ?", creds.Username).Update("refresh_token", refreshToken)
		if result.Error != nil {
			return result.Error
		}

		c.SetCookie(refreshCookie)
	}

	return c.NoContent(http.StatusOK)
}

func (h BiuHandler) LogoutHandler(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		return err
	}

	result := h.DB.Model(&UserCredentials{}).Where("refresh_token = ?", refreshToken.Value).Update("refresh_token", nil)
	if result.Error != nil {
		return result.Error
	}
	c.SetCookie(GenerateCookie(TOKEN_COOKIE_NAME, "", -1))
	c.SetCookie(GenerateCookie(REFRESH_COOKIE_NAME, "", -1))
	return c.NoContent(http.StatusOK)
}

func (h BiuHandler) RefreshHandler(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh")
	if err != nil {
		return err
	}

	var username string
	result := h.DB.Model(&UserCredentials{}).Select("username").First(&username, "refresh_token = ?", refreshToken.Value)
	if result.Error != nil {
		return result.Error
	}

	tokenCookie, err := GenerateTokenCookie(username)
	if err != nil {
		return err
	}
	c.SetCookie(tokenCookie)

	newRefreshToken, refreshCookie, err := GenerateRefreshCookie(username)
	if err != nil {
		return err
	}

	result = h.DB.Model(&UserCredentials{}).Where("username = ?", username).Update("refresh_token", newRefreshToken)
	if result.Error != nil {
		return result.Error
	}

	c.SetCookie(refreshCookie)

	return c.NoContent(http.StatusOK)
}

func (h BiuHandler) ForgotHandler(c echo.Context) error {
	var creds CredentialsRequest
	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		return err
	}

	var userId int
	result := h.DB.Model(&AppUser{}).Select("id").First(&userId, "email = ?", creds.Username)
	if result.Error != nil {
		return result.Error
	}

	randomToken, err := GenerateRandomToken()
	if err != nil {
		return err
	}

	err = h.RC.Set(context.Background(), randomToken, userId, 24*time.Second).Err()
	if err != nil {
		return err
	}

	// send email
	return c.NoContent(http.StatusOK)
}

func (h BiuHandler) RecoverHandler(c echo.Context) error {
	token := c.QueryParam("tk")
	if token == "" {
		return errors.New("no token")
	}

	userId, err := h.RC.GetDel(context.Background(), token).Result()
	if err != nil {
		return err
	}

	fmt.Printf("userId %v", userId)
	return c.NoContent(http.StatusOK)
}

func (h BiuHandler) VerifyHandler(c echo.Context) error {
	token := c.QueryParam("tk")
	if token == "" {
		return errors.New("no token")
	}

	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return err
	}

	var userId int
	result := h.DB.Model(&AppUser{}).Select("id").First(&userId, "email = ?", claims.Username)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(http.StatusOK, IdResource{Id: userId})
}
