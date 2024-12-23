package handlers

import (
	"babel/db"
	"babel/models"
	"html/template"
	"net/http"
	"strings"
)

func GenerateMenuFields(db *db.DB) []models.MenuItem {
	// dsn := "myuser:mypassword@tcp(host.docker.internal:3306)/babel?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var raw_menulist []models.DBMenuItem
	db.Raw(`
		SELECT name, description, GROUP_CONCAT(version ORDER BY ranking) as version FROM (
			SELECT d.name, d.description,
			CONCAT(dh.version_major, '.', dh.version_minor, '.', dh.version_patch) as version ,
			RANK() over (partition by d.name order by dh.version_major DESC, dh.version_minor DESC, dh.version_patch DESC) as ranking
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
			Link:  "/docs/" + item.Name + "/" + latest_version + "/",
		})
		for _, version := range versions {
			v := models.MenuItem{
				Title: version,
				Link:  "/docs/" + item.Name + "/" + version + "/",
			}
			versions_links = append(versions_links, v)
		}

		// now setting up the final menu for the dropdown
		menurow := models.MenuItem{
			Title:    item.Name,
			Link:     "/docs/" + item.Name,
			Children: versions_links,
			MoreInfo: "/info/" + item.Name,
		}

		menulist = append(menulist, menurow)

		// fmt.Printf("%s %s %s\n", item.Name, item.Description, item.Version)
	}

	// for _, item := range menulist {
	// 	fmt.Printf("%s %s %s\n", item.Title, item.Link, item.MoreInfo)
	// }
	return menulist
}

func IndexHandler(db *db.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		data := GenerateMenuFields(db)

		page := template.Must(template.ParseFiles("templates/index.html"))
		page.ExecuteTemplate(res, "index.html", data)
	}
}
