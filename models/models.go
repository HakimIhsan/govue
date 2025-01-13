package models

import (
    "gorm.io/gorm"
)

var DB *gorm.DB

type Book struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title      string    `gorm:"size:255" json:"title"`
    Writer      string    `gorm:"size:255" json:"writer"`
    Description string    `gorm:"type:text" json:"description"`
    Date        MyDate	 `json:"date"`
    Price       int       `json:"price"`
    Photo       string    `gorm:"size:255" json:"photo"`
}