package models

import (
    "gorm.io/gorm"
    "time"
)

var DB *gorm.DB

type Book struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title      string    `gorm:"size:255" json:"title"`
    Writer      string    `gorm:"size:255" json:"writer"`
    Description string    `gorm:"type:text" json:"description"`
    Date        time.Time `json:"date"`
    Price       int       `json:"price"`
    Photo       string    `gorm:"size:255" json:"photo"`
}

type User struct {
    ID          int      `gorm:"primaryKey" json:"id"`
    Username      string    `gorm:"size:255" json:"username"`
    Password      string    `gorm:"size:255" json:"password"`
}

type PersonalToken struct {
    ID          int      `gorm:"primaryKey" json:"id"`
    UserId      int    `json:"user_id"`
    Token      string    `gorm:"size:255" json:"token"`
    ExpiredAt  time.Time    `json:"expired_at"`
}
