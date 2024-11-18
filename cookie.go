package main

import (
	"math"
	"net/http"
)

const (
	TOKEN_COOKIE_NAME   = "token"
	REFRESH_COOKIE_NAME = "refresh"
	cookieExpiration    = 86400 // seconds
)

func GenerateTokenCookie(username string) (*http.Cookie, error) {
	token, err := GenerateToken(username)
	if err != nil {
		return &http.Cookie{}, err
	}

	return GenerateCookie(TOKEN_COOKIE_NAME, token, cookieExpiration), nil
}

func GenerateRefreshCookie(username string) (string, *http.Cookie, error) {
	refreshToken, err := GenerateRandomToken()
	if err != nil {
		return "", &http.Cookie{}, err
	}

	refreshCookie := GenerateCookie(REFRESH_COOKIE_NAME, refreshToken, math.MaxInt32)

	return refreshToken, refreshCookie, nil
}

func GenerateCookie(name string, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
		SameSite: http.SameSiteStrictMode,
	}
}
