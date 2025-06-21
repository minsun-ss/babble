package models

// this is the nextjs endpoint for listing all libraries
type DBIndexMenuItem struct {
	ProjectTeam string `gorm:"column:project_name"`
	Library     string `gorm:"column:name"`
}

type JsonIndexMenuItem struct {
	ProjectTeam string   `json:"project_name"`
	Libraries   []string `json:"libraries"`
}

// this is the nextjs endpoint for listing all versions
type DBLibraryMenuItem struct {
	Library            string `gorm:"column:name"`
	ProjectTeam        string `gorm:"column:project_name"`
	LibraryDescription string `gorm:"column:description"`
	Version            string `gorm:"column:version"`
}

type JsonLibraryMenuItem struct {
	Library            string   `json:"library"`
	ProjectTeam        string   `json:"project_name"`
	LibraryDescription string   `json:"description"`
	Versions           []string `json:"versions"`
}
