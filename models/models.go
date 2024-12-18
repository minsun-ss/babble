package models

// models for setting up the full library page
type LibraryLink struct {
	Version string
	Link    string
}

type LibraryData struct {
	Library     string
	Description string
	Links       []LibraryLink
}

type MenuItem struct {
	Title    string
	Link     string
	Children []MenuItem
	MoreInfo string
}
