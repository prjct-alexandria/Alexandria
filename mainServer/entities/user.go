package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

func (u User) Hash() (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	u.Pwd = string(hash)
	return u, err
}
