package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FolderName string `gorm:"not null"`
	FileName   string `json:"file" gorm:"not null"`
	FilePath   string
	UserID     uint
	Safe       bool
}

type Folder struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	ParentID uint   `json:"parent_id"`
}

type FileDelete struct {
	FileID         uint `json:"file" gorm:"not null"`
	AdminApproval1 uint
	AdminApproval2 uint
	AdminApproval3 uint
	Approved       bool
}
