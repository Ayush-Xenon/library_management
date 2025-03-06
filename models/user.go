package models

// import (
// 	"time"
// )

// type User struct {
// 	ID            uint `binding:"required" gorm:"primaryKey"`
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// 	Name          string    `binding:"required"`
// 	Email         string    `gorm:"unique_index"`
// 	ContactNumber string    `binding:"required"`
// 	Role          string    `binding:"required"`
// 	PasswordHash  string    `binding:"required"`
// 	Libraries     []Library `gorm:"many2many:user_libraries;"`
// }
