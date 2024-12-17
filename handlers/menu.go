package handlers

import (
	"babel/models"
)

func HandleMenuItem() []models.MenuItem {
	menu := []models.MenuItem{
		{
			Title: "Menu",
			Link:  "#",
		},
		{
			Title: "traderpythonlib",
			Link:  "/docs",
			Children: []models.MenuItem{
				{Title: "Latest", Link: "/docs"},
				{Title: "1.29.0", Link: "/products/new"},
				{Title: "1.28.0", Link: "/products/categories"},
			},
			More: "/docs",
		},
		{
			Title: "deskbot",
			Link:  "#",
			Children: []models.MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "3.0.0", Link: "/users/new"},
				{Title: "2.9.0", Link: "/users/groups"},
			},
		},
		{
			Title: "fndmoodeng",
			Link:  "#",
			Children: []models.MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "1.0.0", Link: "/users/new"},
				{Title: "0.9.0", Link: "/users/groups"},
			},
		},
	}

	return menu
}
