package models

import (
	"gorm.io/gorm"
)

type Fileobject struct {
	gorm.Model        // Embedding gorm.Model provides ID, CreatedAt, UpdatedAt, DeletedAt fields
	Typeis     string `gorm:"column:typeis"` // Use appropriate field names and tags
	Name       string `gorm:"column:name"`
	Content    string `gorm:"column:content"`
}

type User struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password;default:''" json:"password,omitempty"` // Change to string
	GoogleID string `gorm:"column:picture" json:"picture,omitempty"`
}
