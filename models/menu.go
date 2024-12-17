package models

type MenuItem struct {
	Title    string
	Link     string
	Children []MenuItem
	More     string
}
