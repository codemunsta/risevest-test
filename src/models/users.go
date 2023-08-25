package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstname" gorm:"text;not null;default:null"`
	LastName  string `json:"lastname" gorm:"text;not null;default:null"`
	Email     string `json:"email" gorm:"text;not null;default:null"`
	Password  string `json:"password"`
}
