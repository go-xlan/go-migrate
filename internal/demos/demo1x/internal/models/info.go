package models

import "gorm.io/gorm"

type InfoV1 struct {
	gorm.Model
	Name string `gorm:"unique;column:name"`
	Cate string `gorm:"column:cate"`
}

func (*InfoV1) TableName() string {
	return "infos"
}

type InfoV2 struct {
	gorm.Model
	Name string `gorm:"unique;column:name"`
	Cate int64  `gorm:"column:cate"`
}

func (*InfoV2) TableName() string {
	return "infos"
}

type InfoV3 struct {
	gorm.Model
	Name string `gorm:"unique;column:name"`
	Cate int8   `gorm:"column:cate"`
}

func (*InfoV3) TableName() string {
	return "infos"
}
