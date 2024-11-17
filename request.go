package main

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

func (request UserRequest) ToAppUser() AppUser {
	creds := UserCredentials{Username: request.Email, Password: request.Password}
	return AppUser{
		Name:            request.Name,
		Email:           request.Email,
		UserCredentials: creds,
	}
}
