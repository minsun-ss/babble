package handlers

import (
	"babel/db"
	"babel/models"
)

type MenuItem struct {
	Title    string
	Link     string
	Children []MenuItem
	MoreInfo string
}

func HandleMenuItem() []models.MenuItem {
	menu := db.GenerateMenuFields()

	return menu
}
