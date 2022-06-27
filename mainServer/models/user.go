package models

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type RegisterForm struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

type LoginForm struct {
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}
