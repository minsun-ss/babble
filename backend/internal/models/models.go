/*
Package models contains all the models (aka structs) required for running
the application, including gorm-specific models.

Prefixes denote if they are models specific to certain packges:
  - DB: Gorm models
  - Page: structs specific for HTML page specific items
*/
package models

import (
	"html/template"
)

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
type PageLibraryLink struct {
	Version string
	Link    string
}

type PageLibraryData struct {
	Library     string
	Description string
	Links       []PageLibraryLink
}

// for the index pages
type PageMenuItem struct {
	Title    string
	Link     string
	Children []PageMenuItem
	MoreInfo string
}

type PageIndex struct {
	MenuItems []PageMenuItem
	Body      template.HTML
}

type ZipResult struct {
	Value []byte
}
