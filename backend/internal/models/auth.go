package models

type DBUsername struct {
	Username string `gorm:"column:username"`
}

func (DBUsername) TableName() string {
	return "users"
}

type DBProjectNameResult struct {
	ProjectName string `gorm:"column:project_name"`
}

type DBUserInsert struct {
	Username string `gorm:"column:username"`
	Role     string `gorm:"column:role"`
	Iat      int64  `gorm:"column:iat"`
}

func (DBUserInsert) TableName() string {
	return "users"
}
