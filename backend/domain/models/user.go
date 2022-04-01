package models

import (
	"time"
)

// Pitch  User Model
type User struct {
	ID *uint64 `json:"id"`

	Mobile *string `json:"mobile" gorm:"type:varchar(10); unique; not null"`
	Name   *string `json:"name" gorm:"type:varchar(50);"`
}

// Pitch Admin User Model
type AdminUser struct {
	ID        *uint64    `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`

	// Firebase UID
	UID   *string `json:"uid" gorm:"type:varchar(50); unique; not null"`
	Email *string `json:"email" gorm:"type:varchar(50); unique; not null"`

	Role *string `json:"role" gorm:"type:varchar(50);"`
}
