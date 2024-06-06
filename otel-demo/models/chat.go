package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Prompt string `gorm:"column:prompt"`
	Reply  string `gorm:"column:reply"`
}
