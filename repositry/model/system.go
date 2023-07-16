package model

import "gorm.io/gorm"

type System struct {
	gorm.Model
	Name string
	Text string
}
