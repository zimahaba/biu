package transport

import "github.com/zimahaba/biu/internal/models"

type CredentialsRequest struct {
	Username     string
	Password     string
	KeepLoggedIn bool
}

type UserRequest struct {
	Name     string
	Email    string
	Password string
}

func (request UserRequest) ToAppUser() models.AppUser {
	creds := models.UserCredentials{Username: request.Email, Password: request.Password}
	return models.AppUser{
		Name:            request.Name,
		Email:           request.Email,
		UserCredentials: creds,
	}
}
