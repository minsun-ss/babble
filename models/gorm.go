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
