package models

import (
	"time"
)

type User struct {
	ID            uint `binding:"required" gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string    `binding:"required"`
	Email         string    `gorm:"unique_index"`
	ContactNumber string    `binding:"required"`
	Role          string    `binding:"required"`
	Password      string    `binding:"required"`
	Libraries     []Library `gorm:"foreignKey:ID;many2many:user_libraries;OnDelete:CASCADE; OnUpdate:CASCADE;"`
}

type Library struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Users []User `gorm:"foreignKey:ID;many2many:user_libraries; OnDelete:CASCADE; OnUpdate:CASCADE;"`
}

type LibraryInput struct {
	Name string `binding:"required"`
}

type UserLibraries struct {
	UserID    uint 
	LibraryID uint 
}

type AuthInput struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}
type AuthCreate struct {
	Email         string `binding:"required"`
	Password      string `binding:"required"`
	Name          string `binding:"required"`
	ContactNumber string `binding:"required"`
	//Role string `binding:"required"`
}

type Book struct {
	ISBN            string `gorm:"primary_key"`
	LibID           uint   `gorm:"primary_key"`
	Title           string
	Authors         string
	Publisher       string
	Version         string
	TotalCopies     int
	AvailableCopies int
}

type BookInput struct {
	ISBN            string `binding:"required"`


}

type RequestEvent struct {
	ReqID        uint `gorm:"primary_key"`
	BookID       string
	ReaderID     uint
	RequestDate  string
	ApprovalDate string
	ApproverID   uint
	RequestType  string
	LibID        uint
}

type IssueRegistry struct {
	IssueID            uint `gorm:"primary_key"`
	ISBN               string
	ReaderID           uint
	IssueApproverID    uint
	IssueStatus        string
	IssueDate          string
	ExpectedReturnDate string
	ReturnDate         string
	ReturnApproverID   uint
}

type UserClaims struct {
	ID uint `json:"id"`
}

type ValidateOutput struct {
	Result  bool
	Message string
}
