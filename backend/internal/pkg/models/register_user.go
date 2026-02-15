package models

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Success bool `json:"success"`
}
