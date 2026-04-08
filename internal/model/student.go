package model

import (
	"time"
)

// Student represents the student entity in the system
type Student struct {
	ID           int64      `json:"id" db:"id"`
	FirstName    string     `json:"first_name" db:"first_name" binding:"required,max=50"`
	LastName     string     `json:"last_name" db:"last_name" binding:"required,max=50"`
	Age          int        `json:"age" db:"age" binding:"required,gt=0,lt=150"`
	Gender       string     `json:"gender" db:"gender" binding:"required,oneof=Male Female Other"`
	Email        string     `json:"email" db:"email" binding:"required,email,max=255"`
	Phone        string     `json:"phone" db:"phone" binding:"required,max=15"`
	Class        string     `json:"class" db:"class" binding:"required,max=20"`
	Rank         *string    `json:"rank" db:"rank" binding:"len=1,oneof=A B C D E F"`
	AddressLine1 string     `json:"address_line1" db:"address_line1" binding:"required,max=100"`
	AddressLine2 string     `json:"address_line2" db:"address_line2" binding:"max=100"`
	City         string     `json:"city" db:"city" binding:"required,max=50"`
	State        string     `json:"state" db:"state" binding:"required,max=50"`
	Pincode      string     `json:"pincode" db:"pincode" binding:"required,len=6,numeric"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}
