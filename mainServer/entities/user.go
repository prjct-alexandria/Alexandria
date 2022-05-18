package entities

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"omitempty"`
}
