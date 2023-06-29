package model

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	gorm.Model
	FileId      string `gorm:"primaryKey"`
	UserId      string `gorm:"not null,index"`
	MD5         string `gorm:"not null,index"`
	ParentId    string `gorm:"index"`
	Size        int64  `gorm:"not null"`
	Name        string `gorm:"not null,index"`
	Cover       string
	Path        string `gorm:"not null"`
	IsDirectory bool   `gorm:"not null"`
	Category    int
	Type        int
	Status      int        `gorm:"index"`
	Flag        int        `gorm:"index"`
	RestoredAt  *time.Time `gorm:"index"`
}
