package models

type DBUserName struct {
	Username string `gorm:"column:username"`
}

type DBProjectName struct {
	ProjectName string `gorm:"column:project_name"`
}
