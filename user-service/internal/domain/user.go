package domain

import (
	"time"
)

type User struct {
	ID        uint    `gorm:"primaryKey"`
	Username  string  `gorm:"uniqueIndex;not null"`
	Email     string  `gorm:"uniqueIndex;not null"`
	Following []*User `gorm:"many2many:user_followers;joinForeignKey:UserID;joinReferences:FollowerID"`
	Followers []*User `gorm:"many2many:user_followers;joinForeignKey:FollowerID;joinReferences:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
