package handlers

import (
	"babel/models"
	"babel/utils"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

// generate the fields for the menu on index.html
func generateMenuFields(db *utils.DB) []models.MenuItem {
	var rawMenuList []models.DBMenuItem

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
		GROUP BY name, description;`).Scan(&rawMenuList)

	var menuList []models.MenuItem
	for _, item := range rawMenuList {
		// setting up the children
		versions := strings.Split(item.Version, ",")
		latestVersion := versions[0]

		var versionsLinks []models.MenuItem
		versionsLinks = append(versionsLinks, models.MenuItem{
			Title: "Latest Version",
			Link:  "/docs/" + item.Name + "/" + latestVersion + "/",
		})
		for _, version := range versions {
			v := models.MenuItem{
				Title: version,
				Link:  "/docs/" + item.Name + "/" + version + "/",
			}
			versionsLinks = append(versionsLinks, v)
		}

		// now setting up the final menu for the dropdown
		menuRow := models.MenuItem{
			Title:    item.Name,
			Link:     "/docs/" + item.Name,
			Children: versionsLinks,
			MoreInfo: "/info/" + item.Name,
		}

		menuList = append(menuList, menuRow)

		slog.Debug("Loaded menu item", "name", item.Name, "description", item.Description,
			"version", item.Version)
	}
	return menuList
}

// handles the "/" endpoint
func IndexHandler(db *utils.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		data := generateMenuFields(db)

		page := template.Must(template.ParseFiles("templates/index.html"))
		page.ExecuteTemplate(res, "index.html", data)
	}
}
