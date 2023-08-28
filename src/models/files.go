package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FolderName string
	FileName   string
	FilePath   string
	UserID     uint
	Safe       bool
}

type Folder struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Name     string `gorm:"not null"`
	ParentID uint
	User     User
}
