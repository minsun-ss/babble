package models

// this is the nextjs endpoint for listing all libraries
type DBIndexMenuItem struct {
	ProjectTeam string `gorm:"column:project_team"`
	Library     string `gorm:"column:name"`
}

type JsonIndexMenuItem struct {
	ProjectTeam string   `json:"project_team"`
	Libraries   []string `json:"libraries"`
}

// this is the nextjs endpoint for listing all versions
type DBLibraryMenuItem struct {
	Library            string `gorm:"column:name"`
	ProjectTeam        string `gorm:"column:project_team"`
	LibraryDescription string `gorm:"column:description"`
	Version            string `gorm:"column:version"`
}

type JsonLibraryMenuItem struct {
	Library            string   `json:"library"`
	ProjectTeam        string   `json:"project_team"`
	LibraryDescription string   `json:"description"`
	Versions           []string `json:"versions"`
}
