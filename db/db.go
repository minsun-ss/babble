package db

import (
	"babel/models"
	"fmt"

	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MenuListResult struct {
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Version     string `gorm:"column:version"`
}

func FetchAllLibraryInfo(path string) models.LibraryData {
	data := models.LibraryData{
		Library:     path,
		Description: "A trading library",
		Links: []models.LibraryLink{
			{Version: "1.2.3", Link: "/docs"},
			{Version: "1.2.1", Link: "/docs"},
		},
	}

	return data
}

func GenerateMenuFields() []models.MenuItem {
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
			MoreInfo: "/info/traderpythonlib",
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

func Stuff() []models.MenuItem {
	dsn := "myuser:mypassword@tcp(host.docker.internal:3306)/babel?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var raw_menulist []MenuListResult
	db.Raw(`
		SELECT name, description, GROUP_CONCAT(version) as version FROM (
			SELECT d.name, d.description,
			CONCAT(dh.version_major, '.', dh.version_minor, '.', dh.version_patch) as version ,
			RANK() over (partition by d.name order by dh.version_major, dh.version_minor, dh.version_patch) as ranking
			FROM babel.docs d
			JOIN babel.doc_history dh
			on d.name=dh.name
			WHERE hidden = 0
			) as versions
		WHERE ranking < 6
		GROUP BY name, description;`).Scan(&raw_menulist)

	var menulist []models.MenuItem
	for _, item := range raw_menulist {
		// setting up the children
		versions := strings.Split(item.Version, ",")
		latest_version := versions[0]

		var versions_links []models.MenuItem
		versions_links = append(versions_links, models.MenuItem{
			Title: "Latest Version",
			Link:  "/docs/" + item.Name + "/" + latest_version,
		})
		for _, version := range versions {
			v := models.MenuItem{
				Title: version,
				Link:  "/docs/" + item.Name + "/" + version,
			}
			versions_links = append(versions_links, v)
		}

		row := models.MenuItem{
			Title:    item.Name,
			Link:     "/docs/" + item.Name,
			Children: versions_links,
			MoreInfo: "/info/" + item.Name,
		}

		menulist = append(menulist, row)

		fmt.Printf("%s %s %s\n", item.Name, item.Description, item.Version)
	}

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	for _, item := range menulist {
		fmt.Printf("%s %s %s\n", item.Title, item.Link, item.MoreInfo)
	}
	return menulist
}
