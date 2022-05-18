package entities

type Article struct {
	Id     string
	Title  string
	Owners []User
}
