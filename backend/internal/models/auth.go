package models

type DBUsername struct {
	Username string `gorm:"column:username"`
}

func (DBUsername) TableName() string {
	return "users"
}

type DBUserKey struct {
	Username string `gorm:"column:username"`
	Role     string `gorm:"column:role"`
	IAT      int64  `gorm:"column:iat"`
}

func (DBUserKey) TableName() string {
	return "users"
}

type DBProjectName struct {
	ProjectName string `gorm:"column:project_name"`
}

func (DBProjectName) TableName() string {
	return "projects"
}

type DBUserAccess struct {
	Username    string `gorm:"column:username"`
	ProjectName string `gorm:"column:project_name"`
}

func (DBUserAccess) TableName() string {
	return "user_access"
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
