package handlers

import ("zippygo/models", "html/template", "net.http")

func handleMenu() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", serve)

	// Menu data structure
	menu := []MenuItem{
		{
			Title: "Menu",
			Link:  "#",
		},
		{
			Title: "traderpythonlib",
			Link:  "/docs",
			Children: []MenuItem{
				{Title: "Latest", Link: "/docs"},
				{Title: "1.29.0", Link: "/products/new"},
				{Title: "1.28.0", Link: "/products/categories"},
			},
			More: "/docs",
		},
		{
			Title: "deskbot",
			Link:  "#",
			Children: []MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "3.0.0", Link: "/users/new"},
				{Title: "2.9.0", Link: "/users/groups"},
			},
		},
		{
			Title: "fndmoodeng",
			Link:  "#",
			Children: []MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "1.0.0", Link: "/users/new"},
				{Title: "0.9.0", Link: "/users/groups"},
			},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", menu)
	})
}
