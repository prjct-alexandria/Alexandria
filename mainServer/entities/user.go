package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	Name  string ``
	Email string
	Pwd   string
}

type StrippedUser struct {
	Name  string
	Email string
}

func (u User) Hash() (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	u.Pwd = string(hash)
	return u, err
}
