package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Name     string
	Gender   string
	Address  string
	Avatar   string
	DOB      time.Time
	Role     string
}
