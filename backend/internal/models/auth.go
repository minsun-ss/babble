package models

type DBUsername struct {
	Username string `gorm:"column:username"`
}

func (DBUsername) TableName() string {
	return "users"
}

type DBProjectName struct {
	ProjectName string `gorm:"column:project_name"`
}

func (DBProjectName) Tablename() string {
	return "projects"
}

type DBUserInsert struct {
	Username string `gorm:"column:username"`
	Role     string `gorm:"column:role"`
	Iat      int64  `gorm:"column:iat"`
}

func (DBUserInsert) TableName() string {
	return "users"
}

type DBProjectInsert struct {
	ProjectName string  `gorm:"column:project_name"`
	Email       *string `gorm:"column:email"`
}

func (DBProjectInsert) TableName() string {
	return "projects"
}
