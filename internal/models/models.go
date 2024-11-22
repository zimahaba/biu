package models

import "database/sql"

type UserCredentials struct {
	ID           int
	Username     string
	Password     string
	RefreshToken sql.NullString
	AppUserId    int
}

type AppUser struct {
	ID              int
	Name            string
	Email           string
	Birthday        sql.NullTime
	UserCredentials UserCredentials
}

func (AppUser) TableName() string {
	return "app_user"
}
