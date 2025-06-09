package example2

import "gorm.io/gorm"

type UserV1 struct {
	gorm.Model
	Username string
}

func (*UserV1) TableName() string {
	return "users"
}

type UserV2 struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Nickname string
}

func (*UserV2) TableName() string {
	return "users"
}
