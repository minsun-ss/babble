/*
Package models contains all the models required for the application,
including gorm-specific models. Gorm specific models are prefixed by the
name DB.
*/
package models

// gorm result for MenuListItem
type DBMenuItem struct {
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Version     string `gorm:"column:version"`
}

type DBLibraryItem struct {
	Description string `gorm:"column:description"`
	Version     string `gorm:"column:version"`
}

type DBLibraryZip struct {
	DataZip []byte `gorm:"column:html"`
}

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

type ZipResult struct {
	Value []byte
}
