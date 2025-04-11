package handlers

import (
	"babel/backend/internal/models"
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// generate the fields for the menu on index.html
func generateMenuFields(db *gorm.DB) []models.PageMenuItem {
	var rawMenuList []models.DBMenuItem

	db.Raw(`
		SELECT name, description, GROUP_CONCAT(version ORDER BY ranking) as version FROM (
			SELECT d.name, d.description,
			CONCAT(dh.version_major, '.', dh.version_minor, '.', dh.version_patch) as version ,
			RANK() over (partition by d.name order by dh.version_major DESC, dh.version_minor DESC, dh.version_patch DESC) as ranking
			FROM babel.docs d
			JOIN babel.doc_history dh
			on d.name=dh.name
			WHERE is_visible = 1
			) as versions
		WHERE ranking < 6
		GROUP BY name, description;`).Scan(&rawMenuList)

	var menuList []models.PageMenuItem
	for _, item := range rawMenuList {
		// setting up the children
		versions := strings.Split(item.Version, ",")
		latestVersion := versions[0]

		var versionsLinks []models.PageMenuItem
		versionsLinks = append(versionsLinks, models.PageMenuItem{
			Title: "Latest Version",
			Link:  "/docs/" + item.Name + "/" + latestVersion + "/",
		})
		for _, version := range versions {
			v := models.PageMenuItem{
				Title: version,
				Link:  "/docs/" + item.Name + "/" + version + "/",
			}
			versionsLinks = append(versionsLinks, v)
		}

		// now setting up the final menu for the dropdown
		menuRow := models.PageMenuItem{
			Title:    item.Name,
			Link:     "/docs/" + item.Name,
			Children: versionsLinks,
			MoreInfo: "/info/" + item.Name,
		}

		menuList = append(menuList, menuRow)

		slog.Debug("loaded menu item", "name", item.Name, "description", item.Description,
			"version", item.Version)
	}
	return menuList
}

// IndexHandler handles the "/" endpoint
func IndexHandler(db *gorm.DB, babelFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// remember that the staticFS doesn't have the same path as the handler
		staticHtml, err := babelFS.ReadFile("assets/indexContent.html")

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		data := generateMenuFields(db)
		pageIndexData := models.PageIndex{
			MenuItems: data,
			Body:      template.HTML(string(staticHtml)),
		}

		page := template.Must(template.ParseFS(babelFS, "assets/templates/index.html"))
		page.ExecuteTemplate(w, "index.html", pageIndexData)
	}
}
