package models

import "gorm.io/gorm"

type UserV1 struct {
	gorm.Model
	Username string `gorm:"unique"`
	Nickname string `gorm:"column:nickname"`
	Rank     string `gorm:"column:rank"`
	Score    string `gorm:"column:score"`
}

func (*UserV1) TableName() string {
	return "users"
}

type UserV2 struct {
	gorm.Model
	Username string `gorm:"unique"`
	Nickname string `gorm:"column:nickname"`
	Rank     uint64 `gorm:"column:rank"`
	Score    string `gorm:"column:score"`
}

func (*UserV2) TableName() string {
	return "users"
}

type UserV3 struct {
	gorm.Model
	Username string  `gorm:"unique"`
	Nickname string  `gorm:"column:nickname"`
	Rank     uint64  `gorm:"column:rank"`
	Score    float64 `gorm:"column:score"`
}

func (*UserV3) TableName() string {
	return "users"
}
