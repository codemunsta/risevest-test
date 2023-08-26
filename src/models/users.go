package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstname" gorm:"text;not null;default:null"`
	LastName  string `json:"lastname" gorm:"text;not null;default:null"`
	Email     string `json:"email" gorm:"text;not null;default:null"`
	Password  string `json:"password" gorm:"text;not null"`
	Files     []File
	Folders   []Folder
}

type Admin struct {
	gorm.Model
	FirstName string `json:"firstname" gorm:"not null"`
	LastName  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	Role      string `json:"-" gorm:"type:enum('super_admin', 'admin')"`
}
